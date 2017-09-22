package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gholt/brimtext"
)

func table(args []string) {
	if len(args) == 0 {
		help("", 0)
	}
	switch args[0] {
	case "search":
		tableSearch(args[1:])
	case "search-column":
		tableSearchColumn(args[1:])
	default:
		help(fmt.Sprintf("Unknown command table %q.", args[0]), 1)
	}
}

func errnil(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func tableSearch(args []string) {
	if len(args) < 2 {
		help("table search needs a <file> and <phrase>", 1)
	}
	header, data := tableRead(args[0])
	phrase := strings.ToLower(strings.Join(args[1:], " "))
	report := [][]string{header, nil}
	for _, row := range data {
		for _, column := range row {
			if strings.Contains(strings.ToLower(column), phrase) {
				report = append(report, row)
			}
		}
	}
	fmt.Print(brimtext.Align(report, brimtext.NewSimpleAlignOptions()))
}

func tableSearchColumn(args []string) {
	if len(args) < 3 {
		help("table search-column needs a <file>, <column>, and <phrase>", 1)
	}
	header, data := tableRead(args[0])
	columnSearch := args[1]
	columnMatch := -1
	for columnIndex, column := range header {
		if strings.ToLower(column) == strings.ToLower(columnSearch) {
			columnMatch = columnIndex
			break
		}
	}
	if columnMatch == -1 {
		errnil(fmt.Errorf("Could not find column %q", columnSearch))
	}
	phrase := strings.ToLower(strings.Join(args[2:], " "))
	report := [][]string{header, nil}
	for _, row := range data {
		if strings.Contains(strings.ToLower(row[columnMatch]), phrase) {
			report = append(report, row)
		}
	}
	fmt.Print(brimtext.Align(report, brimtext.NewSimpleAlignOptions()))
}

func tableRead(filename string) (header []string, data [][]string) {
	f, err := os.Open(filename)
	errnil(err)
	scanner := bufio.NewScanner(f)
	lineNumber := 0
	shouldBeNoMoreLines := false
	for scanner.Scan() {
		if shouldBeNoMoreLines {
			errnil(fmt.Errorf("trailing lines after line number %d", lineNumber))
		}
		lineNumber++
		strs := strings.Split(scanner.Text(), "|")
		if len(strs) == 1 {
			if lineNumber != 1 && lineNumber != 3 {
				shouldBeNoMoreLines = true
			}
			continue
		}
		if len(strs) < 3 {
			errnil(fmt.Errorf("line number %d has too few columns", lineNumber))
		}
		if strs[0] != "" || strs[len(strs)-1] != "" {
			errnil(fmt.Errorf("line number %d is malformed", lineNumber))
		}
		strs = strs[1 : len(strs)-1]
		if len(data) > 0 && len(strs) != len(data[0]) {
			errnil(fmt.Errorf("line number %d has incorrect number of columns; had %d and expected %d", lineNumber, len(strs), len(data[0])))
		}
		for i, s := range strs {
			strs[i] = strings.TrimSpace(s)
		}
		data = append(data, strs)
	}
	errnil(scanner.Err())
	if len(data) < 2 {
		errnil(fmt.Errorf("no data"))
	}
	return data[0], data[1:]
}
