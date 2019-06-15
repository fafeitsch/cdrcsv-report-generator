package cdrcsv

import (
	"reflect"
	"testing"
)

func TestFile_RemoveParallelCalls(t *testing.T) {
	file, err := ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	expectedFile, err := ReadWithoutHeaderFromFile("../mockdata/smallcdr_wo_parallelcalls.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actualFile := file.CloneWithParallelCallsRemoved()
	if len(expectedFile.Records) != len(actualFile.Records) {
		t.Errorf("Expected number of records is %d, but was %d", len(expectedFile.Records), len(actualFile.Records))
	}
	typ := reflect.TypeOf(Record{})
	for index, record := range expectedFile.Records {
		expectedVal := reflect.ValueOf(*record)
		actualVal := reflect.ValueOf(*(actualFile.Records[index]))
		for i := 0; i < typ.NumField(); i++ {
			if expectedVal.Field(i).String() != actualVal.Field(i).String() {
				t.Errorf("Expected value for record at line %d for field %s is \"%s\", but was \"%s\".", index+1, typ.Field(i).Name, expectedVal.Field(i).String(), actualVal.Field(i).String())
			}
		}
	}
}
