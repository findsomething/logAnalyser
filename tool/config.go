package tool

import (
	"os"
	"fmt"
	"path/filepath"
)

type MyConfig struct {
	ReadFile   string
	SeekFile   string
	ResultFile string
	ReadLimit  int
	Seek       int
	WorkPath   string
}

func NewConfig(readFile, seekFileName, resultFileName string, readLimit int) *MyConfig {
	workPath, err := os.Getwd()
	PanicCheck(err)

	seekFile := filepath.Join(workPath, "output", seekFileName)
	resultFile := filepath.Join(workPath, "output", resultFileName)

	return &MyConfig{ReadFile:readFile, SeekFile:seekFile,
		ResultFile:resultFile, ReadLimit:readLimit, WorkPath:workPath}
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