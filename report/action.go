package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
)

func applyCountings(countings []CountingsDefinition, records []*cdrcsv.Record) map[string]int {
	result := make(map[string]int)
	for _, counting := range countings {
		counter := 0
		for _, record := range records {
			if counting.Formula.MatchRecord(record) {
				counter = counter + 1
			}
		}
		result[counting.Name] = counter
	}
	return result
}
