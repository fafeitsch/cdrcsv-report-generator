package report

import (
	"fmt"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	html "html/template"
	"io"
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

	templateDefinition, err := text.ParseFiles(settings.TemplateFile)
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

	templateDefinition, err := html.ParseFiles(settings.TemplateFile)
	if err != nil {
		return fmt.Errorf("could not parse the template file %s: %v", settings.TemplateFile, err)
	}
	generatedReport := Report{Stats: applyCountings(reportDefinition.Countings, file.Records), Records: statsFile}
	return templateDefinition.Execute(settings.Writer, generatedReport)
}
