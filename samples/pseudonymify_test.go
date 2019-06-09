package samples

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"os"
	"testing"
	"time"
)

func TestPseudonymify(t *testing.T) {
	reader, err := os.Open("create_test_cdr.csv")
	if err != nil {
		t.Errorf("could not open the test cdr: %v", err)
		return
	}
	defer func() {
		_ = reader.Close()
	}()
	file, err := cdrcsv.ReadWithoutHeader(reader)
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	pseudoParticipants := []Participant{
		{Name: "Rebecca Fazioli", Extension: "02823623452"},
		{Name: "Gunther Herstein", Extension: "023632462"},
		{Name: "Bianca Nueré", Extension: "58"},
		{Name: "Ted Walter", Extension: "16"},
	}
	pseudoContexts := []string{"context1", "context2"}
	data := PseudoData{Participants: pseudoParticipants, Contexts: pseudoContexts}
	settings := Settings{HideChannels: true, HideAppData: true}
	_ = Pseudoymify(&[]cdrcsv.File{file}, data, settings)
	tmp, _ := os.Create("/tmp/test.csv")
	file.WriteAsCsvWithoutHeader(tmp)
	defer tmp.Close()
}

func TestFindParticipants(t *testing.T) {
	reader, err := os.Open("create_test_cdr.csv")
	if err != nil {
		t.Errorf("could not open the test cdr: %v", err)
		return
	}
	defer func() {
		_ = reader.Close()
	}()
	file, err := cdrcsv.ReadWithoutHeader(reader)
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	participants := findParticipants([]cdrcsv.File{file})
	expectedParticipants := []Participant{
		{Name: "", Extension: "0815"},
		{Name: "", Extension: "0923526333"},
		{Name: "John Doe", Extension: "12"},
		{Name: "Judith Queston", Extension: "15"},
	}
	if len(participants) != len(expectedParticipants) {
		t.Errorf("Expected participants were %d, but actual were %d", len(expectedParticipants), len(participants))
		return
	}
	for index, participant := range participants {
		if participant.Name != expectedParticipants[index].Name {
			t.Errorf("Expected Name of participant %d is %s, but was %s", index, expectedParticipants[index].Name, participant.Name)
		}
		if participant.Extension != expectedParticipants[index].Extension {
			t.Errorf("Expected Extension of participant %d is %s, but was %s", index, expectedParticipants[index].Extension, participant.Extension)
		}
	}
}

func TestFindContexts(t *testing.T) {
	reader, err := os.Open("create_test_cdr.csv")
	if err != nil {
		t.Errorf("could not open the test cdr: %v", err)
		return
	}
	defer func() {
		_ = reader.Close()
	}()
	file, err := cdrcsv.ReadWithoutHeader(reader)
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	contexts := findContexts([]cdrcsv.File{file})
	expectedContexts := []string{"external", "internal"}
	if len(contexts) != len(expectedContexts) {
		t.Errorf("Expected contexts were %d, but actual contexts were %d.", len(expectedContexts), len(contexts))
		return
	}
	for index, context := range contexts {
		if context != expectedContexts[index] {
			t.Errorf("Context at location %d should be %s, but was %s", index, expectedContexts[index], context)
		}
	}
}

func TestShiftTime(t *testing.T) {
	aDateString := "2019-06-09 10:03:23"
	aDate, _ := time.Parse(cdrcsv.DateFormat, aDateString)

	shifter := NaturalTimeShifter{}
	modifiedTime := shifter.shiftTime(aDate)

	if aDateString != modifiedTime.Format(cdrcsv.DateFormat) {
		t.Errorf("With the default time shifter, the date %s should not be changed, but was changed to %s.", aDateString, modifiedTime.Format(cdrcsv.DateFormat))
	}
}
