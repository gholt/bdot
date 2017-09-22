package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gholt/brimtext"
)

func csvToTable(args []string) {
	data, err := csv.NewReader(os.Stdin).ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data = append(data, nil)
	copy(data[2:], data[1:])
	data[1] = nil
	fmt.Print(brimtext.Align(data, brimtext.NewSimpleAlignOptions()))
}
