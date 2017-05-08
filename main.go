package main

import (
	"logAnalyser/tool"
	"github.com/droundy/goopt"
	"fmt"
	"os"
)

var readFile, resultFileName, seekFileName, mode, outputPath *string
var readLimit *int

func init() {
	readFile = goopt.String([]string{"--input"}, "", "nginx access log file(require)")
	outputPath = goopt.String([]string{"--outputPath"}, "", "the path for result(option)")
	resultFileName = goopt.String([]string{"--outputName"}, "", "output file(require)")
	seekFileName = goopt.String([]string{"--tmpFile"}, "", "tmp file(option)")
	readLimit = goopt.Int([]string{"--limit"}, 1000, "read num every times(option default 1000)")
	mode = goopt.String([]string{"--mode"}, "", "run mode: spec: redirect to file tail after analysing nums of data and create tmpFile everyDay (option)")
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

	myConfig := tool.NewConfig(*readFile, *outputPath, *seekFileName, *resultFileName, *mode, *readLimit)

	logAnalyser := tool.NewLogAnalyser(myConfig)
	logAnalyser.Analysis()
}