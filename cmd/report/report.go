package main

import (
	"flag"
	"fmt"
	"github.com/fafeitsch/open-callopticum/report"
	"os"
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s: %s [parameters] cdr_file_name\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	definition := flag.String("definition", "./definition.json", "Path to the json file containing the generatedReport definition.")
	templateFile := flag.String("template", "./template.gohtml", "Path to the html file containing the template.")
	flag.Parse()
	err := report.GenerateReport(report.Settings{Writer: os.Stdout, CdrFile: flag.Arg(0), ReportDefFile: *definition, TemplateFile: *templateFile})
	if err != nil {
		fmt.Printf("error while creating report: %v", err)
		os.Exit(1)
	}
}
