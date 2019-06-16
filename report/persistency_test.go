package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"os"
	"testing"
)

func TestParseDefinition(t *testing.T) {
	reportDefinition, err := ParseDefinitionFromFile("../mockdata/reportDefinition.json")
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(reportDefinition.Countings) != 3 {
		t.Errorf("Expected number of parsed countings should be 2 but was %d.", len(reportDefinition.Countings))
	}
	expectedNames := []string{"production_calls", "evening_hours", "employees"}
	expectedDisplayNames := []string{"Calls from the Production Department", "Calls from the headquarter in the evening hours", "Added calls from Magdalene Greenman and Farlie Brager"}
	for index, counting := range reportDefinition.Countings {
		if expectedNames[index] != counting.Name {
			t.Errorf("Expected name of the %dth counting is '%s' but was '%s.'", index, expectedNames[index], counting.Name)
		}
		if expectedDisplayNames[index] != counting.DisplayName {
			t.Errorf("Expected display name of the %dth counting is '%s' but was '%s'.", index, expectedDisplayNames[index], counting.DisplayName)
		}
	}
}

func TestParseMatcher(t *testing.T) {
	jsonFile, err := os.Open("../mockdata/reportDefinition.json")
	if err != nil {
		t.Errorf("%v", err)
	}
	defer func() {
		_ = jsonFile.Close()
	}()
	reportDefinition, err := ParseDefinition(jsonFile)
	if err != nil {
		t.Errorf("%v", err)
	}
	var records = []cdrcsv.Record{cdrcsv.Record{Dcontext: "production"},
		{Dcontext: "hq", Start: "2017-06-12 17:23:05"},
		{Dcontext: "production", Start: "2017-06-12 17:23:05"},
		{Dcontext: "hq", Start: "2017-06-12 06:17:22"}}

	expectedMatcher1 := []bool{true, false, true, false}
	expectedMatcher2 := []bool{false, true, false, false}
	prodMatcher := reportDefinition.Countings[0].Formula
	hourMatcher := reportDefinition.Countings[1].Formula
	for index, record := range records {
		if prodMatcher.MatchRecord(&record) != expectedMatcher1[index] {
			t.Errorf("On Record %d, the production matcher reported match = %t, but expected was %t.", index, !expectedMatcher1[index], expectedMatcher1[index])
		}
		if hourMatcher.MatchRecord(&record) != expectedMatcher2[index] {
			t.Errorf("On Record %d, the hour matcher reported match = %t, but expected was %t.", index, !expectedMatcher2[index], expectedMatcher2[index])
		}
	}
}
