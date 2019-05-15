package samples

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"strings"
)

func pseudoymify(cdrs []cdrcsv.File, contacts []SampleContact) []cdrcsv.Record {
	return nil
}

func findRealNames(cdrs []cdrcsv.File) map[string]string {
	realNames := make(map[string]string)
	for _, cdr := range cdrs {
		for _, record := range cdr.Records {
			callerId := record.CallerId
			firstBracket := strings.Index(callerId, "<")
			name := callerId[1 : firstBracket-2]
			realNames[name] = callerId[firstBracket+1 : len(callerId)-1]
		}
	}
	return realNames
}
