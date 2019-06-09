package cdrcsv

import (
	"os"
	"testing"
	"time"
)

func TestPseudonymify(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/cdr.csv")
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	personsFile, err := os.Open("../mockdata/persons.csv")
	if err != nil {
		t.Errorf("could not open the persons file: %v", err)
		return
	}
	defer func() {
		_ = personsFile.Close()
	}()
	pseudoParticipants, err := ParsePseudoContacts(personsFile, false)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	pseudoContexts := []string{"c1", "c2", "c3", "c4", "c5"}
	data := PseudoData{Participants: pseudoParticipants, Contexts: pseudoContexts}
	settings := Settings{HideChannels: true, HideAppData: true}
	err = Pseudonymify(&[]File{file}, data, settings)
	if err != nil {
		t.Errorf("could not pseudonymify the dataset: %v", err)
	}
}

func TestPseudonymifyTooFewContexts(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/cdr.csv")
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	personsFile, err := os.Open("../mockdata/persons.csv")
	if err != nil {
		t.Errorf("could not open the persons file: %v", err)
		return
	}
	defer func() {
		_ = personsFile.Close()
	}()
	pseudoParticipants, err := ParsePseudoContacts(personsFile, false)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	pseudoContexts := []string{"c1", "c2", "c3"}
	data := PseudoData{Participants: pseudoParticipants, Contexts: pseudoContexts}
	settings := Settings{HideChannels: true, HideAppData: true}
	err = Pseudonymify(&[]File{file}, data, settings)
	if err != nil {
		return
	}
	t.Errorf("pseudonymify should return an error because to few contexts were given")
}

func TestPseudonymifyTooFewContacts(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/cdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	pseudoParticipants := []Participant{
		{Name: "Rebecca Fazioli", Extension: "02823623452"},
		{Name: "Gunther Herstein", Extension: "023632462"},
		{Name: "Bianca Nueré", Extension: "58"},
		{Name: "Ted Walter", Extension: "16"},
	}
	pseudoContexts := []string{"c1", "c2", "c3", "c4", "c5"}
	data := PseudoData{Participants: pseudoParticipants, Contexts: pseudoContexts}
	settings := Settings{HideChannels: true, HideAppData: true}
	err = Pseudonymify(&[]File{file}, data, settings)
	if err != nil {
		return
	}
	t.Errorf("pseudonymify should return an error because to few contacts were given.")
}

func TestFindParticipants(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/cdr.csv")
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	participants := findParticipants([]File{file})
	expectedParticipants := []Participant{
		{Name: "", Extension: "103-358-0893"},
		{Name: "", Extension: "108-939-4916"},
		{Name: "", Extension: "109-503-4250"},
		{Name: "", Extension: "127-110-7139"},
		{Name: "", Extension: "146-846-6697"},
		{Name: "", Extension: "153-585-7133"},
		{Name: "", Extension: "195-189-2657"},
		{Name: "", Extension: "228-679-2771"},
		{Name: "", Extension: "235-992-7436"},
		{Name: "", Extension: "239-553-6964"},
		{Name: "", Extension: "251-219-5099"},
		{Name: "", Extension: "251-791-1061"},
		{Name: "", Extension: "256-666-8232"},
		{Name: "", Extension: "260-342-4531"},
		{Name: "", Extension: "261-773-6574"},
		{Name: "", Extension: "262-599-6351"},
		{Name: "", Extension: "300-677-9221"},
		{Name: "", Extension: "325-447-5931"},
		{Name: "", Extension: "329-762-8080"},
		{Name: "", Extension: "334-442-8436"},
		{Name: "", Extension: "360-603-8614"},
		{Name: "", Extension: "360-690-5326"},
		{Name: "", Extension: "367-620-8069"},
		{Name: "", Extension: "378-254-9697"},
		{Name: "", Extension: "386-769-6230"},
		{Name: "", Extension: "397-815-2211"},
		{Name: "", Extension: "413-365-0861"},
		{Name: "", Extension: "418-342-5659"},
		{Name: "", Extension: "418-700-7488"},
		{Name: "", Extension: "434-142-9024"},
		{Name: "", Extension: "462-658-8898"},
		{Name: "", Extension: "475-619-9098"},
		{Name: "", Extension: "502-679-5153"},
		{Name: "", Extension: "514-826-1358"},
		{Name: "", Extension: "532-527-4907"},
		{Name: "", Extension: "538-170-4665"},
		{Name: "", Extension: "556-508-6671"},
		{Name: "", Extension: "564-633-4692"},
		{Name: "", Extension: "606-166-6300"},
		{Name: "", Extension: "640-254-7297"},
		{Name: "", Extension: "641-611-9779"},
		{Name: "", Extension: "642-769-3957"},
		{Name: "", Extension: "651-813-5422"},
		{Name: "", Extension: "656-217-5184"},
		{Name: "", Extension: "666-901-7816"},
		{Name: "", Extension: "712-450-5280"},
		{Name: "", Extension: "715-413-9112"},
		{Name: "", Extension: "730-868-5047"},
		{Name: "", Extension: "734-432-8226"},
		{Name: "", Extension: "749-409-5989"},
		{Name: "", Extension: "760-695-0607"},
		{Name: "", Extension: "785-220-1129"},
		{Name: "", Extension: "787-952-0687"},
		{Name: "", Extension: "791-445-9811"},
		{Name: "", Extension: "803-921-5515"},
		{Name: "", Extension: "816-361-1520"},
		{Name: "", Extension: "826-313-0409"},
		{Name: "", Extension: "834-314-7166"},
		{Name: "", Extension: "842-627-1145"},
		{Name: "", Extension: "846-758-4856"},
		{Name: "", Extension: "859-403-3109"},
		{Name: "", Extension: "864-117-6588"},
		{Name: "", Extension: "866-935-7752"},
		{Name: "", Extension: "870-526-4260"},
		{Name: "", Extension: "872-619-0407"},
		{Name: "", Extension: "877-539-8401"},
		{Name: "", Extension: "877-911-5347"},
		{Name: "", Extension: "889-969-2806"},
		{Name: "", Extension: "914-510-3340"},
		{Name: "", Extension: "929-488-7484"},
		{Name: "", Extension: "951-981-7011"},
		{Name: "", Extension: "956-833-3388"},
		{Name: "", Extension: "974-948-4817"},
		{Name: "", Extension: "989-326-7716"},
		{Name: "Bentlee McVicker", Extension: "917-375-0980"},
		{Name: "Farlie Brager", Extension: "490-156-8031"},
		{Name: "Jacobo Lissandri", Extension: "605-142-0198"},
		{Name: "Leanna Cuphus", Extension: "609-326-2780"},
		{Name: "Magdalene Greenman", Extension: "190-590-0260"},
		{Name: "Marget Biernacki", Extension: "253-433-5862"},
		{Name: "Onfre MacFaul", Extension: "992-725-2449"},
		{Name: "Sherilyn Aughton", Extension: "672-769-5651"},
		{Name: "Trace Lavender", Extension: "716-523-6632"},
		{Name: "Tymothy Hamblin", Extension: "748-621-0365"},
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
	file, err := ReadWithoutHeaderFromFile("../mockdata/cdr.csv")
	if err != nil {
		t.Errorf("could not parse test cdr: %v", err)
		return
	}
	contexts := findContexts([]File{file})
	expectedContexts := []string{"door", "hq", "production", "support"}
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
	aDate, _ := time.Parse(DateFormat, aDateString)

	shifter := NaturalTimeShifter{}
	modifiedTime := shifter.shiftTime(aDate)

	if aDateString != modifiedTime.Format(DateFormat) {
		t.Errorf("With the default time shifter, the date %s should not be changed, but was changed to %s.", aDateString, modifiedTime.Format(DateFormat))
	}
}
