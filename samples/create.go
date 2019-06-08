package samples

import (
	"fmt"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"io"
	"math/rand"
	"strconv"
)

type Options struct {
	Count             int
	Contacts          []SampleContact
	Seed              int64
	CompanyExtensions []string
}

type SampleContact struct {
	firstName         string
	lastName          string
	externalExtension string
	internalExtension string
	isEmployee        bool
	internalPhone     string
}

func Create(options *Options, out io.Writer) error {
	if len(options.Contacts) < 2 {
		return fmt.Errorf("extension of contacts is smaller than 2")
	}
	records := make([]*cdrcsv.Record, 0, options.Count)
	random := rand.New(rand.NewSource(options.Seed))
	for i := 0; i < options.Count; i++ {
		callerIndex := random.Intn(len(options.Contacts))
		calleeIndex := random.Intn(len(options.Contacts))
		for callerIndex == calleeIndex {
			calleeIndex = rand.Intn(len(options.Contacts))
		}
		caller := options.Contacts[callerIndex]
		callee := options.Contacts[calleeIndex]
		record := createRecord(caller, callee, random, options)
		records = append(records, &record)
	}
	cdrFile := cdrcsv.File{Records: records}
	return cdrFile.WriteAsCsvWithoutHeader(out)
}

func createRecord(caller SampleContact, callee SampleContact, rnd *rand.Rand, options *Options) cdrcsv.Record {
	record := cdrcsv.Record{}
	record.Accountcode = ""
	if caller.isEmployee {
		record.Src = caller.internalExtension
		record.CallerId = fmt.Sprintf("\"%s %s\" <%s>", caller.firstName, caller.lastName, caller.internalExtension)
		record.Dcontext = "internal"
		record.Channel = fmt.Sprintf("internal-0000%d", rnd.Intn(100))
	} else {
		record.Src = caller.externalExtension
		record.CallerId = "\"\" <" + caller.externalExtension + ">"
		record.Dcontext = "external"
		record.Channel = fmt.Sprintf("external-0000%d", rnd.Intn(100))
	}
	if callee.isEmployee {
		record.Dst = callee.internalExtension
		if !caller.isEmployee {
			record.Dst = options.CompanyExtensions[rnd.Intn(len(options.CompanyExtensions))]
		}
		record.DstChannel = fmt.Sprintf("%s-0000%d", callee.internalPhone, rnd.Intn(100))
		record.LastData = fmt.Sprintf("SIP/%s", callee.internalPhone)
	} else {
		record.Dst = callee.externalExtension
		record.DstChannel = fmt.Sprintf("to_public-0000%d", rnd.Intn(100))
		record.LastData = fmt.Sprintf("SIP/%s@to_public", callee.externalExtension)
	}
	record.LastApp = "DIAL"
	record.Start = "2019-05-09 08:06:11"
	record.Answer = "2019-05-09 08:06:30"
	record.End = "2019-05-09 08:30:12"
	record.Duration = "1441"
	record.Billsec = "1422"
	record.Disposition = cdrcsv.ANSWERED
	record.AmaFlag = cdrcsv.DOCUMENTATION
	record.Userfield = ""
	record.UniqueId = strconv.Itoa(rand.Intn(1250))
	return record
}
