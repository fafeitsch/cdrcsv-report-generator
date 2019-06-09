package samples

import (
	"encoding/csv"
	"fmt"
	"io"
)

func ParseCsv(reader io.Reader, hasHeader bool) ([]Participant, error) {
	csvReader := csv.NewReader(reader)
	lines, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading from csv file: %v", err)
	}
	start := 0
	if hasHeader {
		start = 1
	}
	result := make([]Participant, 0, len(lines))
	for _, line := range lines[start:] {
		contact := Participant{
			Name:      line[0] + " " + line[1],
			Extension: line[2],
		}
		result = append(result, contact)
	}
	return result, nil
}
