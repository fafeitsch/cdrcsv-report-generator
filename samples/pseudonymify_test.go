package samples

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"os"
	"testing"
)

func TestFindParticipants(t *testing.T) {
	reader, err := os.Open("create_test_cdr.csv")
	if err != nil {
		t.Errorf("could not open the test cdr: %v", err)
		return
	}
	file, err := cdrcsv.ReadWithoutHeader(reader)
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	participants := findParticipants([]cdrcsv.File{*file})
	expectedParticipants := []participant{
		{name: "Judith Queston", extension: "15"},
		{name: "John Doe", extension: "12"},
		{name: "", extension: "0815"},
		{name: "", extension: "0923526333"},
	}
	if len(participants) != len(expectedParticipants) {
		t.Errorf("Expected participants were %d, but actual were %d", len(expectedParticipants), len(participants))
		return
	}
	for index, participant := range participants {
		if participant.name != expectedParticipants[index].name {
			t.Errorf("Expected name of participant %d is %s, but was %s", index, expectedParticipants[index].name, participant.name)
		}
		if participant.extension != expectedParticipants[index].extension {
			t.Errorf("Expected extension of participant %d is %s, but was %s", index, expectedParticipants[index].extension, participant.extension)
		}
	}
}
