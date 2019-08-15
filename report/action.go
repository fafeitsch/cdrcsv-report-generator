package report

import (
	"fmt"
	"github.com/fafeitsch/cdrcsv-report-generator/cdrcsv"
	html "html/template"
	"io"
	"path/filepath"
	text "text/template"
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

type Settings struct {
	ReportDefFile  string
	TemplateFile   string
	CdrFile        string
	PlainText      bool
	Writer         io.Writer
	RemoveParallel bool
}

func GeneratePlainTextReport(settings Settings) error {
	file, err := cdrcsv.ReadWithoutHeaderFromFile(settings.CdrFile)
	if err != nil {
		return fmt.Errorf("could not read cdr csv file %s: %v", settings.CdrFile, err)
	}
	if settings.RemoveParallel {
		file = file.CloneWithParallelCallsRemoved()
	}
	statsFile, err := newStatsFile(file)
	if err != nil {
		return fmt.Errorf("could not create statistics: %v", err)
	}
	reportDefinition, err := ParseDefinitionFromFile(settings.ReportDefFile)
	if err != nil {
		return fmt.Errorf("could not parse definition file %s: %v", settings.ReportDefFile, err)
	}

	name := filepath.Base(settings.TemplateFile)
	templateDefinition, err := text.New(name).Funcs(text.FuncMap{
		"diff": func(minuend int, subtrahend int) int {
			return minuend - subtrahend
		},
		"add": func(summand1 int, summand2 int) int {
			return summand1 + summand2
		},
	}).ParseFiles(settings.TemplateFile)
	if err != nil {
		return fmt.Errorf("could not parse the template file %s: %v", settings.TemplateFile, err)
	}
	generatedReport := Report{Stats: applyCountings(reportDefinition.Countings, file.Records), Records: statsFile}
	return templateDefinition.Execute(settings.Writer, generatedReport)
}

func GenerateHtmlReport(settings Settings) error {
	file, err := cdrcsv.ReadWithoutHeaderFromFile(settings.CdrFile)
	if err != nil {
		return fmt.Errorf("could not read cdr csv file %s: %v", settings.CdrFile, err)
	}
	if settings.RemoveParallel {
		file = file.CloneWithParallelCallsRemoved()
	}
	statsFile, err := newStatsFile(file)
	if err != nil {
		return fmt.Errorf("could not create statistics: %v", err)
	}
	reportDefinition, err := ParseDefinitionFromFile(settings.ReportDefFile)
	if err != nil {
		return fmt.Errorf("could not parse definition file %s: %v", settings.ReportDefFile, err)
	}
	name := filepath.Base(settings.TemplateFile)
	templateDefinition, err := html.New(name).Funcs(html.FuncMap{
		"diff": func(minuend int, subtrahend int) int {
			return minuend - subtrahend
		},
		"add": func(summand1 int, summand2 int) int {
			return summand1 + summand2
		},
	}).ParseFiles(settings.TemplateFile)
	if err != nil {
		return fmt.Errorf("could not parse the template file %s: %v", settings.TemplateFile, err)
	}
	generatedReport := Report{Stats: applyCountings(reportDefinition.Countings, file.Records), Records: statsFile}
	return templateDefinition.Execute(settings.Writer, generatedReport)
}
