package main

import (
	"flag"
	"fmt"
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"github.com/fafeitsch/open-callopticum/samples"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	contactsPath := flag.String("contacts", "./contacts.csv", "Path to a file holding the pseudo contacts. Read more about the format in the README.md.")
	contexts := flag.String("contexts", "context1,context2,context3", "Comma-separated list of pseudo contacts.")
	shiftDay := flag.Int("days", 0, "Shift all time relevant data by days. Can be positive or negative")
	shiftYear := flag.Int("years", 0, "Shift all time relevant data by years. Can be positive or negative")
	shiftHour := flag.Int("hours", 0, "Shift all time relevant data by hours. Can be positive or negative")
	shiftMinute := flag.Int("minutes", 0, "Shift all time relevant data by minutes. Can be positive or negative")
	hideChannels := flag.Bool("hideChannels", true, "If true then all channel data will be overriden with a fixed string")
	hideData := flag.Bool("hideData", true, "If true then the last app data column is override with a fixed string")
	prefix := flag.String("prefix", "pseudo_", "Denotes the prefix which is attached to the generated pseudonymified files.")
	flag.Parse()

	fmt.Printf("Reading the pseudo contacts file %s …\n", *contactsPath)
	contactsFile, err := os.Open(*contactsPath)
	if err != nil {
		fmt.Printf("Error: Could not open the pseudo contacts file %s: %v", *contactsPath, err)
		os.Exit(1)
	}
	contacts, err := samples.ParseCsv(contactsFile, false)
	_ = contactsFile.Close()
	if err != nil {
		fmt.Printf("Error: The contacts file %s could not be parsed: %v", *contactsPath, err)
		os.Exit(1)
	}
	fmt.Printf("Reading the pseudo contacts successfull.\n")

	contextsSlice := strings.Split(*contexts, ",")
	data := samples.PseudoData{Participants: contacts, Contexts: contextsSlice}

	shifter := samples.NaturalTimeShifter{Hours: *shiftHour, Days: *shiftDay, Years: *shiftYear, Minutes: *shiftMinute}
	settings := samples.Settings{TimeShifter: &shifter, HideAppData: *hideData, HideChannels: *hideChannels}

	fmt.Printf("Reading the cdr files …\n")
	cdrFiles := make([]cdrcsv.File, 0, len(flag.Args()))
	for _, arg := range flag.Args() {
		reader, err := os.Open(arg)
		if err != nil {
			fmt.Printf("could not open file %s: %v", arg, err)
			os.Exit(1)
		}
		file, err := cdrcsv.ReadWithoutHeader(reader)
		_ = reader.Close()
		if err != nil {
			fmt.Printf("could not parse file %s: %v", arg, err)
			_ = reader.Close()
			os.Exit(1)
		}
		cdrFiles = append(cdrFiles, file)
	}

	fmt.Printf("Found %d cdr files.\n", len(cdrFiles))

	err = samples.Pseudoymify(&cdrFiles, data, settings)
	if err != nil {
		fmt.Printf("Error. Could not pseudonymify the cdrs: %v", err)
		os.Exit(1)
	}

	for index, file := range cdrFiles {
		base := filepath.Base(flag.Arg((index)))
		dir := filepath.Dir(flag.Arg(index))
		output, err := os.Create(dir + "/" + *prefix + base)
		if err != nil {
			fmt.Printf("Error: could open file for writing for argument %s: %v", flag.Arg(index), err)
			_ = output.Close()
			os.Exit(1)
		}
		err = file.WriteAsCsvWithoutHeader(output)
		if err != nil {
			fmt.Printf("Error: could not write out pseudonymified version of %s: %v", flag.Arg(index), err)
		}
		_ = output.Close()
	}

}
