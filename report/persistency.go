package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type jsonFormula struct {
	Operator string       `json:"operator"`
	Column   string       `json:"column"`
	Regex    string       `json:"regex"`
	Left     *jsonFormula `json:"left"`
	Right    *jsonFormula `json:"right"`
}

type jsonCounting struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Formula     jsonFormula `json:"formula"`
}

type jsonReport struct {
	Countings []jsonCounting `json:"countings"`
}

func ParseDefinitionFromFile(filename string) (ReportDefinition, error) {
	jsonFile, err := os.Open("../mockdata/reportDefinition.json")
	if err != nil {
		return ReportDefinition{}, err
	}
	defer func() {
		_ = jsonFile.Close()
	}()
	return ParseDefinition(jsonFile)
}

func ParseDefinition(reader io.Reader) (ReportDefinition, error) {
	var jsonReport jsonReport
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return ReportDefinition{}, fmt.Errorf("could not read from stream: %v", err)
	}
	err = json.Unmarshal(buf.Bytes(), &jsonReport)
	if err != nil {
		return ReportDefinition{}, fmt.Errorf("could not parse json: %v", err)
	}
	return convertJson(jsonReport)
}

func convertJson(report jsonReport) (ReportDefinition, error) {
	countings := make([]CountingsDefinition, 0, len(report.Countings))
	for index, jsonCounting := range report.Countings {
		matcher, err := convertFormula(jsonCounting.Formula)
		if err != nil {
			return ReportDefinition{}, fmt.Errorf("parsing the formula of the %dth counting failed: %v", index, err)
		}
		counting := CountingsDefinition{Name: jsonCounting.Name, DisplayName: jsonCounting.DisplayName, Formula: matcher}
		countings = append(countings, counting)
	}
	return ReportDefinition{Countings: countings}, nil
}

func convertFormula(formula jsonFormula) (Matcher, error) {
	if formula.Left != nil && formula.Right != nil && formula.Operator != "" && formula.Regex == "" && formula.Column == "" {
		leftFormula, err := convertFormula(*formula.Left)
		if err != nil {
			return nil, fmt.Errorf("could not parse left formula: %v", err)
		}
		rightFormula, err := convertFormula(*formula.Right)
		if err != nil {
			return nil, fmt.Errorf("could not parse right formula: %v", err)
		}
		if formula.Operator == "and" || formula.Operator == "&&" {
			return &AndMatcher{left: leftFormula, right: rightFormula}, nil
		}
		if formula.Operator == "or" || formula.Operator == "||" {
			return &OrMatcher{left: leftFormula, right: rightFormula}, nil
		}
		return nil, fmt.Errorf("Operator '%s' not known. Use one of 'and', '&&', 'or' or '||'.", formula.Operator)
	}
	if formula.Left == nil && formula.Right == nil && formula.Operator == "" && formula.Regex != "" && formula.Column != "" {
		refType := reflect.TypeOf(cdrcsv.Record{})
		_, ok := refType.FieldByName(formula.Column)
		if !ok {
			fields := make([]string, 0, refType.NumField())
			for i := 0; i < refType.NumField(); i++ {
				fields = append(fields, refType.Field(i).Name)
			}
			return nil, fmt.Errorf("Column '%s' not found. Column must be one of: %s", formula.Column, strings.Join(fields, ", "))
		}
		provider := func(record *cdrcsv.Record) string {
			val := reflect.ValueOf(*record)
			return val.FieldByName(formula.Column).String()
		}
		expression, err := regexp.Compile(formula.Regex)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid regex: %v", formula.Regex, err)
		}
		return &RegexMatcher{regex: *expression, provider: provider}, nil
	}
	return nil, fmt.Errorf("either left, right and operator fields must be set or regex and column fields must be set.")
}
