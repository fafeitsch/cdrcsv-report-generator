package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type jsonFormula struct {
	Operator string      `json:"operator"`
	Column   string      `json:"column"`
	Regex    string      `json:"regex"`
	Left     jsonFormula `json:"left"`
	Right    jsonFormula `json:"right"`
}

type jsonCounting struct {
	Name               string      `json:"name"`
	DisplayName        string      `json:"display_name"`
	Formula            jsonFormula `json:"formula"`
	ExludeOtherMatches bool        `json:"exclude_other_matches"`
}

type jsonReport struct {
	countings []jsonCounting `json:"countings"`
}

func ParseDefinition(reader *io.Reader) (ReportDefinition, error) {
	var jsonReport jsonReport
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(*reader)
	if err != nil {
		return ReportDefinition{}, fmt.Errorf("could not read from stream: %v", err)
	}
	err = json.Unmarshal(buf.Bytes(), jsonReport)
	if err != nil {
		return ReportDefinition{}, fmt.Errorf("could not parse json: %v", err)
	}
	return ReportDefinition{}, nil
}
