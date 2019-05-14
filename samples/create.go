package samples

import (
	"fmt"
	"io"
	"math/rand"
)

type Options struct {
	Count    int
	Contacts []SampleContact
	Seed     int64
}

type SampleContact struct {
	firstName         string
	lastName          string
	externalExtension string
	internalExtension string
	isEmployee        bool
	internalPhone     string
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

func (c *cdrCsvRecord) toCsvString() string {
	return "\"" + c.accountcode +
		"\",\"" + c.src + "\",\"" + c.dst + "\",\"" + c.dcontext + "\",\"" + c.callerId + "\",\"" + c.channel + "\",\"" + c.dstChannel + "\",\"" + c.lastApp + "\",\"" + c.lastData + "\",\"" + c.start + "\",\"" + c.answer + "\",\"" + c.end + "\"," + c.duration + "," + c.billsec + ",\"" + string(c.disposition) + "\",\"" + string(c.amaFlag) + "\",\"" + c.userfield + "\",\"" + c.uniqueId + "\""
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
	if len(options.Contacts) < 2 {
		return fmt.Errorf("number of contacts is smaller than 2")
	}
	random := rand.New(rand.NewSource(options.Seed))
	for i := 0; i < options.Count; i++ {
		callerIndex := random.Intn(len(options.Contacts))
		calleeIndex := random.Intn(len(options.Contacts))
		for callerIndex == calleeIndex {
			calleeIndex = rand.Intn(len(options.Contacts))
		}
		caller := options.Contacts[callerIndex]
		callee := options.Contacts[calleeIndex]
		record := createRecord(caller, callee, random)
		_, err := io.WriteString(out, record.toCsvString()+"\n")
		if err != nil {
			return fmt.Errorf("could not export record %v: %v", record.toCsvString(), err)
		}
	}
	return nil
}

func createRecord(caller SampleContact, callee SampleContact, rnd *rand.Rand) cdrCsvRecord {
	record := cdrCsvRecord{}
	record.accountcode = ""
	if caller.isEmployee {
		record.src = caller.internalExtension
		record.callerId = fmt.Sprintf("%s %s <%s>", caller.firstName, caller.lastName, caller.internalExtension)
		record.dcontext = "internal"
		record.channel = fmt.Sprintf("internal-0000%d", rnd.Intn(100))
	} else {
		record.src = caller.externalExtension
		record.callerId = caller.externalExtension
		record.dcontext = "external"
		record.channel = fmt.Sprintf("external-0000%d", rnd.Intn(100))
	}
	if callee.isEmployee {
		record.dst = callee.internalExtension
		record.dstChannel = fmt.Sprintf("%s-0000%d", callee.internalPhone, rnd.Intn(100))
		record.lastData = fmt.Sprintf("SIP/%s", callee.internalPhone)
	} else {
		record.dst = callee.externalExtension
		record.dstChannel = fmt.Sprintf("to_public-0000%d", rnd.Intn(100))
		record.lastData = fmt.Sprintf("SIP/%s@to_public", callee.externalExtension)
	}
	record.start = "2019-05-09 08:06:11"
	record.answer = "2019-05-09 08:06:30"
	record.end = "2019-05-09 08:30:12"
	record.duration = "1441"
	record.billsec = "1422"
	record.disposition = ANSWERED
	record.amaFlag = DOCUMENTATION
	record.userfield = ""
	record.uniqueId = string(rnd.Int63())
	return record
}
