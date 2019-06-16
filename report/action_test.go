package report

import (
	"bytes"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"testing"
)

func TestApplyCountings(t *testing.T) {
	def, err := ParseDefinitionFromFile("../mockdata/reportDefinition.json")
	file, err := cdrcsv.ReadWithoutHeaderFromFile("../mockdata/smallcdr.csv")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expectedMap := make(map[string]int)
	expectedMap["production_calls"] = 17
	expectedMap["employees"] = 4
	expectedMap["evening_hours"] = 3
	actualMap := applyCountings(def.Countings, file.Records)
	if len(expectedMap) != len(actualMap) {
		t.Errorf("Expected fields were %d, but actual fields were %d.", len(expectedMap), len(actualMap))
		return
	}
	for key, value := range expectedMap {
		if actualMap[key] != value {
			t.Errorf("Expected value for field \"%s\" was %d, but actual was %d.", key, value, actualMap[key])
		}
	}
}

func TestGenerateHtmlReport(t *testing.T) {
	writer := new(bytes.Buffer)
	settings := Settings{
		Writer:        writer,
		ReportDefFile: "../mockdata/reportDefinition.json",
		TemplateFile:  "../mockdata/reportTemplate.tmpl",
		CdrFile:       "../mockdata/smallcdr.csv",
	}
	err := GenerateHtmlReport(settings)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := writer.String()
	expected := `This is a sample call detail report:

Production Calls: 17
Calls in the evening hours: 3
Calls from Magdalene Greenman and Farlie Brager: 4

The average calling time was approximately 5 minutes.
The median calling time was approximately 2 minutes.

The longest call lasted approximately 19 minutes and happened between &#34;&#34; &lt;334-442-8436&gt; and 397-815-2211.

There were 45 calls in total.`
	if actual != expected {
		t.Errorf("Expected text:\n%s\n\nActual text:\n%s", expected, actual)
	}
}

func TestGenerateHtmlReport_RemoveParallel(t *testing.T) {
	writer := new(bytes.Buffer)
	settings := Settings{
		Writer:         writer,
		ReportDefFile:  "../mockdata/reportDefinition.json",
		TemplateFile:   "../mockdata/reportTemplate.tmpl",
		CdrFile:        "../mockdata/smallcdr.csv",
		RemoveParallel: true,
	}
	err := GenerateHtmlReport(settings)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := writer.String()
	expected := `This is a sample call detail report:

Production Calls: 17
Calls in the evening hours: 2
Calls from Magdalene Greenman and Farlie Brager: 4

The average calling time was approximately 5 minutes.
The median calling time was approximately 2 minutes.

The longest call lasted approximately 19 minutes and happened between &#34;&#34; &lt;334-442-8436&gt; and 397-815-2211.

There were 34 calls in total.`
	if actual != expected {
		t.Errorf("Expected text:\n%s\n\nActual text:\n%s", expected, actual)
	}
}

func TestGeneratePlainReport(t *testing.T) {
	writer := new(bytes.Buffer)
	settings := Settings{
		Writer:        writer,
		ReportDefFile: "../mockdata/reportDefinition.json",
		TemplateFile:  "../mockdata/reportTemplate.tmpl",
		CdrFile:       "../mockdata/smallcdr.csv",
	}
	err := GeneratePlainTextReport(settings)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := writer.String()
	expected := `This is a sample call detail report:

Production Calls: 17
Calls in the evening hours: 3
Calls from Magdalene Greenman and Farlie Brager: 4

The average calling time was approximately 5 minutes.
The median calling time was approximately 2 minutes.

The longest call lasted approximately 19 minutes and happened between "" <334-442-8436> and 397-815-2211.

There were 45 calls in total.`
	if actual != expected {
		t.Errorf("Expected text:\n%s\n\nActual text:\n%s", expected, actual)
	}
}

func TestGeneratePlainReport_RemoveParallel(t *testing.T) {
	writer := new(bytes.Buffer)
	settings := Settings{
		Writer:         writer,
		ReportDefFile:  "../mockdata/reportDefinition.json",
		TemplateFile:   "../mockdata/reportTemplate.tmpl",
		CdrFile:        "../mockdata/smallcdr.csv",
		RemoveParallel: true,
	}
	err := GeneratePlainTextReport(settings)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	actual := writer.String()
	expected := `This is a sample call detail report:

Production Calls: 17
Calls in the evening hours: 2
Calls from Magdalene Greenman and Farlie Brager: 4

The average calling time was approximately 5 minutes.
The median calling time was approximately 2 minutes.

The longest call lasted approximately 19 minutes and happened between "" <334-442-8436> and 397-815-2211.

There were 34 calls in total.`
	if actual != expected {
		t.Errorf("Expected text:\n%s\n\nActual text:\n%s", expected, actual)
	}
}
