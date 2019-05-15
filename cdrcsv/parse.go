package cdrcsv

import (
	"fmt"
	"io"
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
	return "\"" + c.Accountcode +
		"\",\"" + c.Src + "\",\"" + c.Dst + "\",\"" + c.Dcontext + "\",\"" + c.CallerId + "\",\"" + c.Channel + "\",\"" + c.DstChannel + "\",\"" + c.LastApp + "\",\"" + c.LastData + "\",\"" + c.Start + "\",\"" + c.Answer + "\",\"" + c.End + "\"," + c.Duration + "," + c.Billsec + ",\"" + string(c.Disposition) + "\",\"" + string(c.AmaFlag) + "\",\"" + c.Userfield + "\",\"" + c.UniqueId + "\""
}

type File struct {
	Records []Record
}

func (f *File) WriteAsCsvWithoutHeader(writer io.Writer) error {
	for _, record := range f.Records {
		_, err := io.WriteString(writer, record.ToCsvString()+"\n")
		if err != nil {
			return fmt.Errorf("could not export record %v: %v", record.ToCsvString(), err)
		}
	}
	return nil
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
