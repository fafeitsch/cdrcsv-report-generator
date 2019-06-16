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
	plainText := flag.Bool("plain", false, "If true, then special characters are not escaped")
	removeParallel := flag.Bool("removeParallel", true, "If true, Cdr-Records of parallel calls are normalized to one call.")
	flag.Parse()
	var err error
	settings := report.Settings{Writer: os.Stdout, CdrFile: flag.Arg(0), ReportDefFile: *definition, TemplateFile: *templateFile, RemoveParallel: *removeParallel}
	if *plainText {
		err = report.GeneratePlainTextReport(settings)
	} else {
		err = report.GenerateHtmlReport(settings)
	}
	if err != nil {
		fmt.Printf("error while creating report: %v", err)
		os.Exit(1)
	}
}
