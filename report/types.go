package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"regexp"
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

type Report struct {
	Stats              map[string]int
	AverageCallingTime time.Duration
	MedianCallingTime  time.Duration
	NumberOfCalls      int
	LongestCall        cdrcsv.Record
}
