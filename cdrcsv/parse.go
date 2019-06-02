package cdrcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
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
	Duration    string
	Billsec     string
	Disposition CallState
	AmaFlag     AmaFlag
	Userfield   string
	UniqueId    string
}

func (c *Record) ToCsvString() string {
	//TODO: escape all fields, not only the caller id
	escapedCallerId := strings.ReplaceAll(c.CallerId, "\"", "\"\"")
	return "\"" + c.Accountcode +
		"\",\"" + c.Src + "\",\"" + c.Dst + "\",\"" + c.Dcontext + "\",\"" + escapedCallerId + "\",\"" + c.Channel + "\",\"" + c.DstChannel + "\",\"" + c.LastApp + "\",\"" + c.LastData + "\",\"" + c.Start + "\",\"" + c.Answer + "\",\"" + c.End + "\"," + c.Duration + "," + c.Billsec + ",\"" + string(c.Disposition) + "\",\"" + string(c.AmaFlag) + "\",\"" + c.Userfield + "\",\"" + c.UniqueId + "\""
}

type File struct {
	Records []Record
}

//WriteAsCsvWithoutHeader writes a CDR file in CSV format to the specified writer.
func (f *File) WriteAsCsvWithoutHeader(writer io.Writer) error {
	for _, record := range f.Records {
		line := record.ToCsvString()
		_, err := io.WriteString(writer, line+"\n")
		if err != nil {
			return fmt.Errorf("could not export record %v: %v", record.ToCsvString(), err)
		}
	}
	return nil
}

//ReadWithoutHeader reads a CDR csv-file from the specified reader and returns the record file.
func ReadWithoutHeader(reader io.Reader) (*File, error) {
	csvReader := csv.NewReader(reader)
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read csv records: %v", err)
	}
	result := make([]Record, 0, len(lines))
	for _, line := range lines {
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
		record.Duration = line[12]
		record.Billsec = line[13]
		record.Disposition = CallState(line[14])
		record.AmaFlag = AmaFlag(line[15])
		record.Userfield = line[16]
		record.UniqueId = line[17]
		result = append(result, record)
	}
	return &File{Records: result}, nil
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
