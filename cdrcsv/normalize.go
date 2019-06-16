package cdrcsv

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

//RemoveParallelCalls removes all parallel calls from the cdr file. Parallel calls are created in an CDR if
//the Dial-App contains an '&'-sign. Example DIAL(SIP/phone1&SIP/phone2) will create two CDR lines with
//nearly identical content, even the unique id is identical. Be careful to distinguish parallel calls from
//transferred calls, which also are nearly identical but their Start time differs. See also:
//https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+CDR+Specification#Asterisk12CDRSpecification
func (f *File) CloneWithParallelCallsRemoved() File {
	candiates := make(map[int]*[]*Record)
	origLocation := make(map[string]int)
	for index, record := range f.Records {
		callId, _ := strconv.Atoi(strings.Split(record.UniqueId, ".")[1])
		origLocation[record.UniqueId] = index
		if _, ok := candiates[callId]; !ok {
			list := make([]*Record, 0, 0)
			candiates[callId] = &list
			list = append(list, record)
		} else {
			list := candiates[callId]
			peer := (*list)[0]
			if recordsAreParallel(record, peer) {
				appendList := append(*list, record)
				candiates[callId] = &appendList
			}
		}
	}
	result := make([]*Record, 0, len(f.Records))
	for _, records := range candiates {
		result = append(result, getAnsweredRecordIfExists(*records))
	}
	sort.Slice(result, func(i int, j int) bool {
		return origLocation[result[i].UniqueId] < origLocation[result[j].UniqueId]
	})
	return File{Records: result}
}

func recordsAreParallel(r1 *Record, r2 *Record) bool {
	uniqueId1 := strings.Split(r1.UniqueId, ".")
	uniqueId2 := strings.Split(r2.UniqueId, ".")
	time1, _ := strconv.Atoi(uniqueId1[0])
	time2, _ := strconv.Atoi(uniqueId2[0])

	if uniqueId1[1] != uniqueId2[1] {
		return false
	}

	if math.Abs(float64(time1-time2)) > 2 {
		return false
	}

	start1, _ := time.Parse(DateFormat, r1.Start)
	start2, _ := time.Parse(DateFormat, r2.Start)
	diff := start2.Sub(start1).Seconds()
	if math.Abs(float64(diff)) > 2 {
		//If this if-block is executed, the two calls have the same uniqueId[0] and the same sequence number (uniqueId[1])
		//But obviously, their start times differ, so they are transferred calls.
		return false
	}
	return true
}

func getAnsweredRecordIfExists(records []*Record) *Record {
	result := records[0]
	for _, record := range records {
		if record.Disposition == ANSWERED {
			return record
		}
		if record.Disposition == BUSY {
			result = record
		}
	}
	return result
}
