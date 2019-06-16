package cdrcsv

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

func TestReadWrite(t *testing.T) {
	csvFile, err := os.Open("../mockdata/cdr.csv")
	defer csvFile.Close()
	if err != nil {
		t.Errorf("error opening the test cdr file: %v", err)
		return
	}
	file, err := ReadWithoutHeader(csvFile)
	if err != nil {
		t.Errorf("the csv file could not be parsed: %v", err)
		return
	}
	writer := bytes.Buffer{}
	err = file.WriteAsCsvWithoutHeader(&writer)
	csvFile.Seek(0, io.SeekStart)
	expectedReader := bufio.NewReader(csvFile)
	actualReader := bufio.NewReader(&writer)
	for {
		expectedLine, err1 := expectedReader.ReadString('\n')
		actualLine, err2 := actualReader.ReadString('\n')
		if err1 != io.EOF && err2 == io.EOF {
			t.Errorf("actual file is shorter than expected file")
			break
		} else if err1 == io.EOF && err2 != io.EOF {
			t.Errorf("actual file is longer than expected file")
			break
		} else if err1 == io.EOF && err2 == io.EOF {
			break
		} else if err1 != nil || err2 != nil {
			t.Errorf("error while comparing files: %v %v", err1, err2)
		}
		if expectedLine != actualLine {
			t.Errorf("expected vs. actual line:\n%s%s", expectedLine, actualLine)
		}
	}
}

func TestFile_ComputeAverageCallingTime(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := file.ComputeAverageCallingTime()
	expected := 284900 * time.Millisecond
	expected = expected.Round(time.Second)
	if actual != expected {
		t.Errorf("Expected average calling time is %f, but was %f", expected.Seconds(), actual.Seconds())
	}
}

func TestFile_ComputeMeanCallingTime(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := file.ComputeMedianCallingTime()
	expected := 125 * time.Second
	if actual != expected {
		t.Errorf("Expected median calling time is %f, but was %f", expected.Seconds(), actual.Seconds())
	}
}

func TestFile_GetLongestCall(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := file.GetLongestCall()
	expectedUniqueId := "1498559959.14"
	if actual.UniqueId != expectedUniqueId {
		t.Errorf("Expected unique id of longest call is %s, but was %s.", expectedUniqueId, actual.UniqueId)
	}
}
