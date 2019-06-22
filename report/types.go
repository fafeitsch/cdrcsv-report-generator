package report

import (
	"fmt"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"math"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Matcher interface {
	MatchRecord(*cdrcsv.Record) bool
}

type RegexMatcher struct {
	regex    regexp.Regexp
	provider func(*cdrcsv.Record) string
}

func (r *RegexMatcher) MatchRecord(record *cdrcsv.Record) bool {
	text := r.provider(record)
	return r.regex.MatchString(text)
}

type AndMatcher struct {
	left  Matcher
	right Matcher
}

func (a *AndMatcher) MatchRecord(record *cdrcsv.Record) bool {
	return a.left.MatchRecord(record) && a.right.MatchRecord(record)
}

type OrMatcher struct {
	left  Matcher
	right Matcher
}

func (a *OrMatcher) MatchRecord(record *cdrcsv.Record) bool {
	return a.left.MatchRecord(record) || a.right.MatchRecord(record)
}

type CountingsDefinition struct {
	Name        string
	DisplayName string
	Formula     Matcher
}

type ReportDefinition struct {
	Countings []CountingsDefinition
}

type Record struct {
	Accountcode string
	Src         string
	Dst         string
	Dcontext    string
	CallerId    string
	Channel     string
	DstChannel  string
	LastApp     string
	LastData    string
	Start       time.Time
	Answer      time.Time
	End         time.Time
	Duration    time.Duration
	Billsec     time.Duration
	Disposition cdrcsv.CallState
	AmaFlag     cdrcsv.AmaFlag
	Userfield   string
	UniqueId    string
}

func NewRecord(r *cdrcsv.Record) (Record, error) {
	var startTime time.Time
	var err error
	if r.Start != "" {
		startTime, err = time.Parse(cdrcsv.DateFormat, r.Start)
		if err != nil {
			return Record{}, err
		}
	}
	var answerTime time.Time
	if r.Answer != "" {
		answerTime, err = time.Parse(cdrcsv.DateFormat, r.Answer)
		if err != nil {
			return Record{}, err
		}
	}
	var endTime time.Time
	if r.End != "" {
		endTime, err = time.Parse(cdrcsv.DateFormat, r.End)
		if err != nil {
			return Record{}, err
		}
	}
	durationSec, err := strconv.Atoi(r.Duration)
	if err != nil {
		return Record{}, fmt.Errorf("could not parse duration \"%s\": %v", r.Duration, err)
	}
	billSec, err := strconv.Atoi(r.Billsec)
	if err != nil {
		return Record{}, fmt.Errorf("could not parse billsec \"%s\": %v", r.Billsec, err)
	}
	return Record{Accountcode: r.Accountcode, Src: r.Src, Dst: r.Dst, CallerId: r.CallerId, Channel: r.Channel, DstChannel: r.DstChannel, LastApp: r.LastApp, LastData: r.LastData, Start: startTime, Answer: answerTime, End: endTime,
		Duration: time.Duration(durationSec) * time.Second, Billsec: time.Duration(billSec) * time.Second, Disposition: r.Disposition, AmaFlag: r.AmaFlag, Userfield: r.Userfield, UniqueId: r.UniqueId}, nil
}

type StatsFile []*Record

func (s StatsFile) ComputeAverageCallingTime() time.Duration {
	sum := 0.0
	counter := 0.0
	for _, record := range s {
		if record.Billsec == 0 {
			continue
		}
		sum = sum + record.Billsec.Seconds()
		counter = counter + 1
	}
	result := time.Duration(sum*1000/counter) * time.Millisecond
	return result.Round(time.Second)
}

func (s StatsFile) ComputeMedianCallingTime() time.Duration {
	callTimes := make([]float64, 0, len(s))
	for _, record := range s {
		if record.Billsec == 0 {
			continue
		}
		callTimes = append(callTimes, record.Billsec.Seconds())
	}
	if len(callTimes) == 0 {
		return 0
	}
	sort.Float64s(callTimes)
	half := int(math.Ceil(0.5*float64(len(callTimes))) - 1)
	return time.Duration(callTimes[half]) * time.Second
}

func (s StatsFile) GetLongestCall() *Record {
	if len(s) == 0 {
		return nil
	}
	call := s[0]
	for _, record := range s {
		if record.Billsec > call.Billsec {
			call = record
		}
	}
	return call
}

func newStatsFile(file cdrcsv.File) (StatsFile, error) {
	result := make(StatsFile, 0, 0)
	for index, record := range file.Records {
		r, err := NewRecord(record)
		if err != nil {
			return nil, fmt.Errorf("cannot parse record at index %d: %v", index, err)
		}
		result = append(result, &r)
	}
	return result, nil
}

type Report struct {
	Stats   map[string]int
	Records StatsFile
}
