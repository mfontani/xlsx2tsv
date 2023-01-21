package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Version contains the binary version. This is added at build time.
var Version = "uncommitted"

func showHelp(reason string) {
	if len(reason) > 0 {
		fmt.Printf("Error: %s\n", reason)
	}
	fmt.Println("Usage: xlsx2tsv FILE.xlsx [SHEET_NUMBER]")
	fmt.Println("")
	fmt.Println("Converts the first sheet of the given FILE.xslx to TSV")
	fmt.Println("If SHEET_NUMBER is given (1-N), it converts that sheet instead of the first.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  xlsx2tsv foo.xlsx")
	fmt.Println("  xlsx2tsv foo.xlsx 2")
	fmt.Printf("\nThis is xlsx2tsv %s\n", Version)
	if len(reason) > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	for _, v := range os.Args[1:] {
		if v == "-version" || v == "--version" {
			fmt.Printf("%s\n", Version)
			os.Exit(0)
		} else if v == "-help" || v == "--help" {
			showHelp("")
		}
	}
	if len(os.Args) == 1 {
		showHelp("Too few arguments.")
	}
	if len(os.Args) > 3 {
		showHelp("Too many arguments.")
	}
	fileName := os.Args[1]
	sheetNumber := 0
	if len(os.Args) == 3 {
		var err error
		sheetNumber, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad sheet number %s. Need a numeric value. See --help.\n", os.Args[2])
			os.Exit(1)
		}
		if sheetNumber < 1 {
			fmt.Fprintf(os.Stderr, "Bad sheet number %s. Need a positive numeric value. See --help.\n", os.Args[2])
			os.Exit(1)
		}
		sheetNumber = sheetNumber - 1 // 0-indexed
	}
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %s: %s\n", fileName, err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing %s: %s\n", fileName, err)
			os.Exit(1)
		}
	}()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		fmt.Fprintf(os.Stderr, "There seem to be no sheets in %s!\n", fileName)
		os.Exit(1)
	}
	if len(sheets) < sheetNumber+1 {
		fmt.Fprintf(os.Stderr, "Not enough sheets in %s to output sheet %d (have %d sheets)!\n", fileName, sheetNumber+1, len(sheets))
		os.Exit(1)
	}
	rows, err := f.GetRows(sheets[sheetNumber])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get rows from wanted sheet %d in %s: %s\n", sheetNumber+1, fileName, err)
		os.Exit(1)
	}
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = '\t'
	for _, row := range rows {
		err = writer.Write(row)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write TSV: %s\n", err)
			os.Exit(1)
		}
	}
	writer.Flush()
	os.Exit(0)
}
