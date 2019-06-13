package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"os"
	"testing"
)

func TestParseDefinition(t *testing.T) {
	jsonFile, err := os.Open("../mockdata/report.json")
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
	if len(reportDefinition.Countings) != 2 {
		t.Errorf("Expected number of parsed countings should be 2 but was %d.", len(reportDefinition.Countings))
	}
	expectedNames := []string{"production_calls", "after_business"}
	expectedDisplayNames := []string{"Calls from the Production Department", "Calls from the headquarter after business hours"}
	expectedExclude := []bool{false, true}
	for index, counting := range reportDefinition.Countings {
		if expectedNames[index] != counting.Name {
			t.Errorf("Expected name of the %dth counting is '%s' but was '%s.'", index, expectedNames[index], counting.Name)
		}
		if expectedDisplayNames[index] != counting.DisplayName {
			t.Errorf("Expected display name of the %dth counting is '%s' but was '%s'.", index, expectedDisplayNames[index], counting.DisplayName)
		}
		if expectedExclude[index] != counting.ExcludeOtherMatches {
			t.Errorf("Expected 'exclude_other_matchings' of %dth counting is %t but was %t.", index, expectedExclude[index], counting.ExcludeOtherMatches)
		}
	}
}

func TestParseMatcher(t *testing.T) {
	jsonFile, err := os.Open("../mockdata/report.json")
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
		{Dcontext: "hq", Start: "2017-06-12 20:23:05"},
		{Dcontext: "production", Start: "2017-06-12 20:23:05"},
		{Dcontext: "hq", Start: "2017-06-12 06:17:22"}}

	expectedMatcher1 := []bool{true, false, true, false}
	expectedMatcher2 := []bool{false, true, false, false}
	prodMatcher := reportDefinition.Countings[0].Formula
	hourMatcher := reportDefinition.Countings[1].Formula
	for index, record := range records {
		if prodMatcher.MatchRecord(record) != expectedMatcher1[index] {
			t.Errorf("On record %d, the production matcher reported match = %t, but expected was %t.", index, !expectedMatcher1[index], expectedMatcher1[index])
		}
		if hourMatcher.MatchRecord(record) != expectedMatcher2[index] {
			t.Errorf("On record %d, the hour matcher reported match = %t, but expected was %t.", index, !expectedMatcher2[index], expectedMatcher2[index])
		}
	}
}
