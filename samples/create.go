package samples

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Options struct {
	Count int
}

type cdrCsvRecord struct {
	accountcode string
	src         string
	dst         string
	dcontext    string
	callerId    string
	channel     string
	dstChannel  string
	lastApp     string
	lastData    string
	start       string
	answer      string
	end         string
	duration    string
	billsec     string
	disposition callState
	amaFlag     amaFlag
	userfield   string
	uniqueId    string
}

func (c *cdrCsvRecord) toSlice() []string {
	return []string{
		"\"" + c.accountcode + "\"",
		"\"" + c.src + "\"",
		"\"" + c.dst + "\"",
		"\"" + c.dcontext + "\"",
		"\"" + c.callerId + "\"",
		"\"" + c.channel + "\"",
		"\"" + c.dstChannel + "\"",
		"\"" + c.lastApp + "\"",
		"\"" + c.lastData + "\"",
		"\"" + c.start + "\"",
		"\"" + c.end + "\"",
		"\"" + c.duration + "\"",
		"\"" + c.billsec + "\"",
		"\"" + string(c.disposition) + "\"",
		"\"" + string(c.amaFlag) + "\"",
		"\"" + c.userfield + "\"",
		"\"" + c.uniqueId + "\"",
	}
}

type callState string

const (
	ANSWERED  callState = "ANSWERED"
	NO_ANSWER callState = "NO_ANSWER"
	BUSY      callState = "BUSY"
	FAILED    callState = "FAILED"
	UNKNOWN   callState = "UNKNOWN"
)

type amaFlag string

const (
	OMIT          amaFlag = "OMIT"
	BILLING       amaFlag = "BILLING"
	DOCUMENTATION amaFlag = "DOCUMENTATION"
	UNKOWN        amaFlag = "Unkown"
)

func Create(options *Options, out io.Writer) error {
	csvWriter := csv.NewWriter(out)
	for i := 0; i < options.Count; i++ {
		record := cdrCsvRecord{
			accountcode: "",
			src:         "4711",
			dst:         "0815",
			dcontext:    "from_public",
			callerId:    "\"John Doe\" <4711>",
			channel:     "SIP/from_public-0000012",
			dstChannel:  "SIP/deskphone_of_boss_000015a",
			lastApp:     "DIAL",
			lastData:    "SIP/deskphone_of_boss",
			start:       "2019-05-09 08:06:11",
			answer:      "2019-05-09 08:06:30",
			end:         "2019-05-09 08:30:12",
			duration:    "1441",
			billsec:     "1422",
			disposition: ANSWERED,
			amaFlag:     DOCUMENTATION,
			userfield:   "",
			uniqueId:    "2265436.50",
		}
		err := csvWriter.Write(record.toSlice())
		if err != nil {
			return fmt.Errorf("could not export record %v: %v", record.toSlice(), err)
		}
	}
	csvWriter.Flush()
	return nil
}
