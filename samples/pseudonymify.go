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

type PseudoData struct {
	Participants []Participant
	Contexts     []string
}

type Settings struct {
	TimeShifter  TimeShifter
	HideAppData  bool
	HideChannels bool
}

func pseudoymify(cdrs *[]cdrcsv.File, pseudo PseudoData, settings Settings) error {
	participants := findParticipants(*cdrs)
	if len(pseudo.Participants) < len(participants) {
		return errors.New(fmt.Sprintf("number of pseudo contacts is not sufficient, at least %d are needed, only %d were provided", len(participants), len(pseudo.Participants)))
	}
	contexts := findContexts(*cdrs)
	if len(pseudo.Contexts) < len(contexts) {
		return errors.New(fmt.Sprintf("number of pseudo contexts is not sufficient, at least %d are needed, only %d were provided", len(contexts), len(pseudo.Contexts)))
	}
	participantMapping := make(map[Participant]Participant)
	for index, participant := range participants {
		participantMapping[participant] = pseudo.Participants[index]
	}
	contextMapping := make(map[string]string)
	for index, context := range contexts {
		contextMapping[context] = pseudo.Contexts[index]
	}
	if settings.TimeShifter == (TimeShifter)(nil) {
		settings.TimeShifter = &IdentityTimeShifter{}
	}
	for _, file := range *cdrs {
		for _, record := range file.Records {
			caller := callerIdToParticipant(record.CallerId)
			pseudoCaller := participantMapping[caller]
			record.CallerId = pseudoCaller.toCallerId()
			srcParticipant, _ := findParticipantByExtension(participants, record.Src)
			pseudoSrc := participantMapping[srcParticipant]
			record.Src = pseudoSrc.Extension

			dstParticipant, _ := findParticipantByExtension(participants, record.Dst)
			pseudoDst := participantMapping[dstParticipant]
			record.Dst = pseudoDst.Extension

			record.Dcontext = contextMapping[record.Dcontext]

			if settings.HideAppData {
				record.LastData = "NOT_AVAILABLE"
			}
			if settings.HideChannels {
				record.Channel = "NOT_AVAILABLE"
				record.DstChannel = "NOT_AVAILABLE"
			}

			start := record.Start
			if start != "" {
				startTime, err := time.Parse(cdrcsv.DateFormat, start)
				if err != nil {
					return fmt.Errorf("starttime %s could not be parsed: %v", start, err)
				}
				record.Start = settings.TimeShifter.shiftTime(startTime).Format(cdrcsv.DateFormat)
			}

			answered := record.Answer
			if answered != "" {
				answeredTime, err := time.Parse(cdrcsv.DateFormat, answered)
				if err != nil {
					return fmt.Errorf("answered time %s could not be parsed: %v", answered, err)
				}
				record.Answer = settings.TimeShifter.shiftTime(answeredTime).Format(cdrcsv.DateFormat)
			}

			end := record.End
			if end != "" {
				endTime, err := time.Parse(cdrcsv.DateFormat, end)
				if err != nil {
					return fmt.Errorf("end time %s could not be parsed: %v", end, err)
				}
				record.End = settings.TimeShifter.shiftTime(endTime).Format(cdrcsv.DateFormat)
			}
		}
	}
	return nil
}

func callerIdToParticipant(callerId string) Participant {
	firstBracket := strings.Index(callerId, "<")
	extension := callerId[firstBracket+1 : len(callerId)-1]
	name := callerId[1 : firstBracket-2]
	return Participant{Name: name, Extension: extension}
}

func findParticipantByExtension(participants []Participant, extension string) (Participant, error) {
	for _, participant := range participants {
		if participant.Extension == extension {
			return participant, nil
		}
	}
	return Participant{}, errors.New(fmt.Sprintf("could not find a participant with extension %s", extension))
}

type Participant struct {
	Name      string
	Extension string
}

func (p *Participant) toCallerId() string {
	return fmt.Sprintf("\"%s\"<%s>", p.Name, p.Extension)
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
	sort.Slice(result, func(i int, j int) bool {
		if result[i].Name == result[j].Name {
			return result[i].Extension < result[j].Extension
		}
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
