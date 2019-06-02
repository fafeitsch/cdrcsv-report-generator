package samples

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"strings"
)

func pseudoymify(cdrs []cdrcsv.File, contacts []SampleContact) []cdrcsv.Record {
	return nil
}

type participant struct {
	name      string
	extension string
}

func findParticipants(cdrs []cdrcsv.File) []participant {
	//TODO: Document assumption that an extension is bound to at most one name
	result := make([]participant, 0)
	distinctExtensions := make(map[string]bool)
	distinctNames := make(map[string]bool)
	for _, cdr := range cdrs {
		for _, record := range cdr.Records {
			callerId := record.CallerId
			firstBracket := strings.Index(callerId, "<")
			extension := callerId[firstBracket+1 : len(callerId)-1]
			name := callerId[1 : firstBracket-2]
			if _, ok := distinctExtensions[extension]; ok {
				continue
			}
			if _, ok := distinctNames[name]; ok {
				continue
			}
			distinctExtensions[extension] = true
			distinctNames[name] = true
			result = append(result, participant{name: name, extension: extension})
		}
		for _, record := range cdr.Records {
			if _, ok := distinctExtensions[record.Src]; ok {
				continue
			}
			distinctExtensions[record.Src] = true
			result = append(result, participant{name: "", extension: record.Src})
		}
		for _, record := range cdr.Records {
			if _, ok := distinctExtensions[record.Dst]; ok {
				continue
			}
			distinctExtensions[record.Dst] = true
			result = append(result, participant{name: "", extension: record.Dst})
		}
	}
	return result
}
