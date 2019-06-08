package samples

import (
	"errors"
	"fmt"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"sort"
	"strings"
	"time"
)

type TimeShifter interface {
	shiftTime(time.Time) time.Time
}

type IdentityTimeShifter struct {
}

func (i *IdentityTimeShifter) shiftTime(time time.Time) time.Time {
	return time
}

func pseudoymify(cdrs []cdrcsv.File, pseudoContacts []Participant, pseudoContexts []string) error {
	participants := findParticipants(cdrs)
	if len(pseudoContacts) < len(participants) {
		return errors.New(fmt.Sprintf("number of pseudo contacts is not sufficient, at least %d are needed, only %d were provided", len(participants), len(pseudoContacts)))
	}
	contexts := findContexts(cdrs)
	if len(pseudoContexts) < len(contexts) {
		return errors.New(fmt.Sprintf("number of pseudo contexts is not sufficient, at least %d are needed, only %d were provided", len(contexts), len(pseudoContexts)))
	}
	participantMapping := make(map[Participant]Participant)
	for index, participant := range participants {
		participantMapping[participant] = pseudoContacts[index]
	}
	contextMapping := make(map[string]string)
	for index, context := range contexts {
		contextMapping[context] = pseudoContexts[index]
	}
	for _, file := range cdrs {
		for _, record := range file.Records {
			_ = callerIdToParticipant(record.CallerId)
		}
	}
}

func callerIdToParticipant(callerId string) Participant {
	firstBracket := strings.Index(callerId, "<")
	extension := callerId[firstBracket+1 : len(callerId)-1]
	name := callerId[1 : firstBracket-2]
	return Participant{Name: name, Extension: extension}
}

type Participant struct {
	Name      string
	Extension string
}

func findParticipants(cdrs []cdrcsv.File) []Participant {
	//TODO: Document assumption that an Extension is bound to at most one Name
	result := make([]Participant, 0)
	distinctExtensions := make(map[string]bool)
	distinctNames := make(map[string]bool)
	for _, cdr := range cdrs {
		for _, record := range cdr.Records {
			participant := callerIdToParticipant(record.CallerId)
			if _, ok := distinctExtensions[participant.Extension]; ok {
				continue
			}
			if _, ok := distinctNames[participant.Name]; ok {
				continue
			}
			distinctExtensions[participant.Extension] = true
			distinctNames[participant.Name] = true
			result = append(result, participant)
		}
		for _, record := range cdr.Records {
			if _, ok := distinctExtensions[record.Src]; ok {
				continue
			}
			distinctExtensions[record.Src] = true
			result = append(result, Participant{Name: "", Extension: record.Src})
		}
		for _, record := range cdr.Records {
			if _, ok := distinctExtensions[record.Dst]; ok {
				continue
			}
			distinctExtensions[record.Dst] = true
			result = append(result, Participant{Name: "", Extension: record.Dst})
		}
	}
	sort.Slice(result, func(i, j) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func findContexts(cdrs []cdrcsv.File) []string {
	result := make([]string, 0)
	set := make(map[string]bool)
	for _, file := range cdrs {
		for _, record := range file.Records {
			set[record.Dcontext] = true
		}
	}
	for key, value := range set {
		if value {
			result = append(result, key)
		}
	}
	sort.Strings(result)
	return result
}
