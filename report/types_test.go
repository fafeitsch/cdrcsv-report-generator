package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"testing"
	"time"
)

func TestFile_ComputeAverageCallingTime(t *testing.T) {
	file, err := cdrcsv.ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	statsFile, err := newStatsFile(file)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := statsFile.ComputeAverageCallingTime()
	expected := 284900 * time.Millisecond
	expected = expected.Round(time.Second)
	if actual != expected {
		t.Errorf("Expected average calling time is %f, but was %f", expected.Seconds(), actual.Seconds())
	}
}

func TestFile_ComputeMeanCallingTime(t *testing.T) {
	file, err := cdrcsv.ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	statsFile, err := newStatsFile(file)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := statsFile.ComputeMedianCallingTime()
	expected := 125 * time.Second
	if actual != expected {
		t.Errorf("Expected median calling time is %f, but was %f", expected.Seconds(), actual.Seconds())
	}
}

func TestFile_ComputeEmptyStats(t *testing.T) {
	records := make([]*cdrcsv.Record, 0)
	file := cdrcsv.File{Records: records}
	statsFile, err := newStatsFile(file)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := statsFile.ComputeMedianCallingTime()
	if actual != 0 {
		t.Errorf("Median of no calls should be 0, but was %d", actual)
	}
	actual = statsFile.ComputeAverageCallingTime()
	if actual != 0 {
		t.Errorf("Average of no calls should be 0, but was %d.", actual)
	}
}

func TestFile_GetLongestCall(t *testing.T) {
	file, err := cdrcsv.ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	statsFile, err := newStatsFile(file)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := statsFile.GetLongestCall()
	expectedUniqueId := "1498559959.14"
	if actual.UniqueId != expectedUniqueId {
		t.Errorf("Expected unique id of longest call is %s, but was %s.", expectedUniqueId, actual.UniqueId)
	}
}
