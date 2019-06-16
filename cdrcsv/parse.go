package cdrcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02 15:04:05"
)

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
	Start       string
	Answer      string
	End         string
	Duration    time.Duration
	Billsec     time.Duration
	Disposition CallState
	AmaFlag     AmaFlag
	Userfield   string
	UniqueId    string
}

func (c *Record) ToCsvString() string {
	//TODO: escape all fields, not only the caller id
	escapedCallerId := strings.ReplaceAll(c.CallerId, "\"", "\"\"")
	duration := strconv.Itoa(int(c.Duration.Seconds()))
	billsec := strconv.Itoa(int(c.Billsec.Seconds()))
	return "\"" + c.Accountcode +
		"\",\"" + c.Src + "\",\"" + c.Dst + "\",\"" + c.Dcontext + "\",\"" + escapedCallerId + "\",\"" + c.Channel + "\",\"" + c.DstChannel + "\",\"" + c.LastApp + "\",\"" + c.LastData + "\",\"" + c.Start + "\",\"" + c.Answer + "\",\"" + c.End + "\"," + duration + "," + billsec + ",\"" + string(c.Disposition) + "\",\"" + string(c.AmaFlag) + "\",\"" + c.UniqueId + "\",\"" + c.Userfield + "\""
}

type File struct {
	Records []*Record
}

//WriteAsCsvWithoutHeader writes a CDR file in CSV format to the specified writer.
func (f *File) WriteAsCsvWithoutHeader(writer io.Writer) error {
	for _, record := range f.Records {
		line := record.ToCsvString()
		_, err := io.WriteString(writer, line+"\n")
		if err != nil {
			return fmt.Errorf("could not export record: %v", err)
		}
	}
	return nil
}

func (f *File) ComputeAverageCallingTime() time.Duration {
	sum := 0.0
	counter := 0.0
	for _, record := range f.Records {
		if record.Billsec == 0 {
			continue
		}
		sum = sum + record.Billsec.Seconds()
		counter = counter + 1
	}
	result := time.Duration(sum*1000/counter) * time.Millisecond
	return result.Round(time.Second)
}

func (f *File) ComputeMedianCallingTime() time.Duration {
	callTimes := make([]float64, 0, len(f.Records))
	for _, record := range f.Records {
		if record.Billsec == 0 {
			continue
		}
		callTimes = append(callTimes, record.Billsec.Seconds())
	}
	sort.Float64s(callTimes)
	half := int(math.Ceil(0.5*float64(len(callTimes))) - 1)
	return time.Duration(callTimes[half]) * time.Second
}

func (f *File) GetLongestCall() *Record {
	if len(f.Records) == 0 {
		return nil
	}
	call := f.Records[0]
	for _, record := range f.Records {
		if record.Billsec > call.Billsec {
			call = record
		}
	}
	return call
}

func ReadWithoutHeaderFromFile(filename string) (File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return File{}, err
	}
	defer func() {
		_ = file.Close()
	}()
	return ReadWithoutHeader(file)
}

//ReadWithoutHeader reads a CDR csv-file from the specified reader and returns the record file.
func ReadWithoutHeader(reader io.Reader) (File, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = 18
	lines, err := csvReader.ReadAll()
	if err != nil {
		return File{}, fmt.Errorf("could not read csv records: %v", err)
	}
	result := make([]*Record, 0, len(lines))
	for index, line := range lines {
		duration, err := strconv.Atoi(line[12])
		if err != nil {
			return File{}, fmt.Errorf("duration of record in line %d could not be parsed: %v", index+1, err)
		}
		billsec, err := strconv.Atoi(line[13])
		if err != nil {
			return File{}, fmt.Errorf("billset of record in line %d could not be parsed: %v", index+1, err)
		}
		record := Record{}
		record.Accountcode = line[0]
		record.Src = line[1]
		record.Dst = line[2]
		record.Dcontext = line[3]
		record.CallerId = line[4]
		record.Channel = line[5]
		record.DstChannel = line[6]
		record.LastApp = line[7]
		record.LastData = line[8]
		record.Start = line[9]
		record.Answer = line[10]
		record.End = line[11]
		record.Duration = time.Duration(duration) * time.Second
		record.Billsec = time.Duration(billsec) * time.Second
		record.Disposition = CallState(line[14])
		record.AmaFlag = AmaFlag(line[15])
		record.UniqueId = line[16]
		record.Userfield = line[17]
		result = append(result, &record)
	}
	return File{Records: result}, nil
}

type CallState string

const (
	ANSWERED  CallState = "ANSWERED"
	NO_ANSWER CallState = "NO_ANSWER"
	BUSY      CallState = "BUSY"
	FAILED    CallState = "FAILED"
	UNKNOWN   CallState = "UNKNOWN"
)

type AmaFlag string

const (
	OMIT          AmaFlag = "OMIT"
	BILLING       AmaFlag = "BILLING"
	DOCUMENTATION AmaFlag = "DOCUMENTATION"
	UNKOWN        AmaFlag = "Unkown"
)
