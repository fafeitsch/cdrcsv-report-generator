package cdrcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TimeShifter interface {
	shiftTime(time.Time) time.Time
}

type NaturalTimeShifter struct {
	Years   int
	Days    int
	Hours   int
	Minutes int
}

func (n *NaturalTimeShifter) shiftTime(t time.Time) time.Time {
	return t.AddDate(n.Years, 0, n.Days).Add(time.Hour*time.Duration(n.Hours) + time.Minute*time.Duration(n.Minutes))
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

func Pseudonymify(cdrs *[]File, pseudo PseudoData, settings Settings) error {
	participants := findParticipants(*cdrs)
	if len(pseudo.Participants) < len(participants) {
		return errors.New(fmt.Sprintf("number of pseudo contacts is not sufficient, at least %d are needed, only %d were provided", len(participants), len(pseudo.Participants)))
	}
	contexts := findContexts(*cdrs)
	if len(pseudo.Contexts) < len(contexts) {
		return errors.New(fmt.Sprintf("number of pseudo contexts is not sufficient, at least %d are needed, only %d were provided", len(contexts), len(pseudo.Contexts)))
	}
	participantMapping := make(map[Participant]*Participant)
	for index, participant := range participants {
		participantMapping[participant] = &pseudo.Participants[index]
		if participant.Name == "" {
			participantMapping[participant].Name = ""
		}
	}
	contextMapping := make(map[string]string)
	for index, context := range contexts {
		contextMapping[context] = pseudo.Contexts[index]
	}
	if settings.TimeShifter == (TimeShifter)(nil) {
		settings.TimeShifter = &NaturalTimeShifter{}
	}
	for fileIndex, file := range *cdrs {
		for recordIndex, record := range file.Records {
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
				startTime, err := time.Parse(DateFormat, start)
				if err != nil {
					return fmt.Errorf("file %d, line %d: %v", fileIndex+1, recordIndex+1, err)
				}
				newTime := settings.TimeShifter.shiftTime(startTime)
				record.Start = newTime.Format(DateFormat)
				epoch := strconv.Itoa(int(newTime.Unix()))
				callId := strings.Split(record.UniqueId, ".")[1]
				record.UniqueId = epoch + "." + callId
			}

			answered := record.Answer
			if answered != "" {
				answeredTime, err := time.Parse(DateFormat, answered)
				if err != nil {
					return fmt.Errorf("file %d, line %d: %v", fileIndex+1, recordIndex+1, err)
				}
				record.Answer = settings.TimeShifter.shiftTime(answeredTime).Format(DateFormat)
			}

			end := record.End
			if end != "" {
				endTime, err := time.Parse(DateFormat, end)
				if err != nil {
					return fmt.Errorf("file %d, line %d: %v", fileIndex+1, recordIndex+1, err)
				}
				record.End = settings.TimeShifter.shiftTime(endTime).Format(DateFormat)
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
	return fmt.Sprintf("\"%s\" <%s>", p.Name, p.Extension)
}

func findParticipants(cdrs []File) []Participant {
	result := make([]Participant, 0)
	distinctExtensions := make(map[string]bool)
	for _, cdr := range cdrs {
		participants := make(map[Participant]bool)
		for _, record := range cdr.Records {
			participant := callerIdToParticipant(record.CallerId)
			if _, ok := participants[participant]; ok {
				continue
			}
			participants[participant] = true
			distinctExtensions[participant.Extension] = true
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

func findContexts(cdrs []File) []string {
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

func ParsePseudoContacts(reader io.Reader, hasHeader bool) ([]Participant, error) {
	csvReader := csv.NewReader(reader)
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading from csv file: %v", err)
	}
	start := 0
	if hasHeader {
		start = 1
	}
	result := make([]Participant, 0, len(lines))
	for _, line := range lines[start:] {
		contact := Participant{
			Name:      line[0] + " " + line[1],
			Extension: line[2],
		}
		result = append(result, contact)
	}
	return result, nil
}
