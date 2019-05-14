package samples

import (
	"encoding/csv"
	"fmt"
	"io"
)

func ParseCsv(reader io.Reader, hasHeader bool) ([]SampleContact, error) {
	csvReader := csv.NewReader(reader)
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading from csv file: %v", err)
	}
	start := 0
	if hasHeader {
		start = 1
	}
	result := make([]SampleContact, 0, len(lines))
	for _, line := range lines[start:] {
		isEmployee := false
		if line[4] == "true" {
			isEmployee = true
		}
		contact := SampleContact{
			firstName:         line[0],
			lastName:          line[1],
			externalExtension: line[2],
			internalExtension: line[3],
			isEmployee:        isEmployee,
			internalPhone:     line[5],
		}
		result = append(result, contact)
	}
	return result, nil
}
