package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"os"
	"testing"
)

func TestApplyCountings(t *testing.T) {
	jsonFile, err := os.Open("../mockdata/reportDefinition.json")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer func() {
		_ = jsonFile.Close()
	}()
	def, err := ParseDefinition(jsonFile)
	file, err := cdrcsv.ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expectedMap := make(map[string]int)
	expectedMap["production_calls"] = 17
	expectedMap["employees"] = 4
	expectedMap["evening_hours"] = 3
	actualMap := applyCountings(def.Countings, file.Records)
	if len(expectedMap) != len(actualMap) {
		t.Errorf("Expected fields were %d, but actual fields were %d.", len(expectedMap), len(actualMap))
		return
	}
	for key, value := range expectedMap {
		if actualMap[key] != value {
			t.Errorf("Expected value for field \"%s\" was %d, but actual was %d.", key, value, actualMap[key])
		}
	}

}
