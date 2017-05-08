package tool

import (
	"os"
	"fmt"
	"path/filepath"
	"time"
)

type MyConfig struct {
	ReadFile   string
	SeekFile   string
	ResultFile string
	ReadLimit  int
	Seek       int
	OutputPath string
	Mode       string
}

func NewConfig(readFile, outputPath, seekFileName, resultFileName, mode string, readLimit int) *MyConfig {
	var err error

	if outputPath == "" {
		outputPath, err = os.Getwd()
		PanicCheck(err)
		outputPath = filepath.Join("output")
	}

	if mode == "spec" {
		seekFileName = fmt.Sprintf("%s.%s", seekFileName, time.Now().Format("2006-01-02"))
	}

	seekFile := filepath.Join(outputPath, seekFileName)
	resultFile := filepath.Join(outputPath, resultFileName)

	return &MyConfig{ReadFile:readFile, SeekFile:seekFile,
		ResultFile:resultFile, ReadLimit:readLimit, OutputPath:outputPath, Mode:mode}
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func PanicCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func WarnCheck(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}