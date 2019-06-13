package report

import (
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
	expectedDisplayNames := []string{"Calls from the Production Department", "Calls after business hours"}
	expectedExclude := []bool{true, false}
	for index, counting := range reportDefinition.Countings {
		if expectedNames[index] != counting.Name {
			t.Errorf("Expected name of the %d.th counting is %s but was %s.", index, expectedNames[index], counting.Name)
		}
		if expectedDisplayNames[index] != counting.DisplayName {
			t.Errorf("Expected display name of the %d.th counting is %s but was %s.", index, expectedDisplayNames[index], counting.DisplayName)
		}
		if expectedExclude[index] != counting.ExcludeOtherMatches {
			t.Errorf("Expected 'exclude_other_matchings' of %d.th counting is %t but was %t.", index, expectedExclude[index], counting.ExcludeOtherMatches)
		}
	}
}
