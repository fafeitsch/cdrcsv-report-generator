package cdrcsv

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
)

func TestReadWrite(t *testing.T) {
	csvFile, err := os.Open("../samples/create_test_cdr.csv")
	defer csvFile.Close()
	if err != nil {
		t.Errorf("error opening the test cdr file: %v", err)
	}
	file, err := ReadWithoutHeader(csvFile)
	if err != nil {
		t.Errorf("the csv file could not be parsed: %v", err)
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
		} else if err1 == io.EOF && err2 != io.EOF {
			t.Errorf("actual file is longer than expected file")
		} else if err1 == io.EOF && err2 == io.EOF {
			break
		} else if err1 != nil || err2 != nil {
			t.Errorf("error while comparing files: %v %v", err1, err2)
		}
		if expectedLine != actualLine {
			t.Errorf("expected vs. actual line:\n%s\n%s", expectedLine, actualLine)
		}
	}
}
