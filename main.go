package main

import (
	"logAnalyser/tool"
	"github.com/droundy/goopt"
	"fmt"
	"os"
)

var readFile, resultFileName, seekFileName *string
var readLimit *int

func init() {
	readFile = goopt.String([]string{"--input"}, "", "nginx access log file")
	resultFileName = goopt.String([]string{"--output"}, "", "output file")
	seekFileName = goopt.String([]string{"--tmpFile"}, "", "tmp file")
	readLimit = goopt.Int([]string{"--limit"}, 1000, "read num every times")
}

func badUsage() {
	fmt.Println(goopt.Usage())
	os.Exit(1)
}

func main() {
	goopt.Parse(nil)

	if *readFile == "" || *resultFileName == "" {
		badUsage()
	}

	if *seekFileName == "" {
		*seekFileName = fmt.Sprintf("%s.tmp", *resultFileName)
	}

	myConfig := tool.NewConfig(*readFile, *seekFileName, *resultFileName, *readLimit)

	logAnalyser := tool.NewLogAnalyser(myConfig)
	logAnalyser.Analysis()
}