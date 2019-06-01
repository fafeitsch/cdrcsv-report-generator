package samples

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	johnDoe := SampleContact{
		firstName:         "John",
		lastName:          "Doe",
		externalExtension: "4711",
		internalExtension: "12",
		isEmployee:        true,
		internalPhone:     "PHONE_1",
	}
	janeFox := SampleContact{
		firstName:         "Jane",
		lastName:          "Fox",
		externalExtension: "0815",
		internalExtension: "14",
		isEmployee:        false,
		internalPhone:     "PHONE_2",
	}
	judithQuestion := SampleContact{
		firstName:         "Judith",
		lastName:          "Queston",
		externalExtension: "2356",
		internalExtension: "15",
		isEmployee:        true,
		internalPhone:     "PHONE_3",
	}
	options := Options{
		Count:             10,
		Contacts:          []SampleContact{johnDoe, janeFox, judithQuestion},
		Seed:              23,
		CompanyExtensions: []string{"0923526332", "0923526333"},
	}
	writer := bytes.Buffer{}
	err := Create(&options, &writer)
	if err != nil {
		t.Errorf("Creating threw unexpected error: %v", err)
	}
	actualReader := bufio.NewReader(&writer)
	fileReader, _ := os.Open("create_test_cdr.csv")
	defer func() {
		_ = fileReader.Close()
	}()
	expectedReader := bufio.NewReader(fileReader)
	for i := 0; i < options.Count; i++ {
		actualLine, _ := actualReader.ReadString('\n')
		expectedLine, _ := expectedReader.ReadString('\n')
		if actualLine == "" {
			t.Errorf("Line %d is empty", i)
		} else if actualLine != expectedLine {
			t.Errorf("Line %d differs (expected vs. actual):\n%s%s", i, expectedLine, actualLine)
		}
	}
}
