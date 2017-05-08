package tool

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"time"
	"math"
	"encoding/json"
	"strconv"
	"strings"
)

var err error

type LogAnalyser struct {
	seek            int64
	myConfig        *MyConfig
	logHandler      *LogHandler
	analysisResults []*LogLineAnalysisResult
}

type LogAnalysisResult struct {
	UpdatedTime             string `json:"update_time"`
	Count2xx                int `json:"2xx_count"`
	Count3xx                int `json:"3xx_count"`
	Count4xx                int `json:"4xx_count"`
	Count5xx                int `json:"5xx_count"`
	RequestAvgTime          float64 `json:"request_avg_time"`
	UpstreamConnectAvgTime  float64 `json:"upstream_connect_avg_time"`
	UpstreamHeaderAvgTime   float64 `json:"upstream_header_avg_time"`
	UpstreamResponseAvgTime float64 `json:"upstream_response_avg_time"`
	RequestMaxTime          float64 `json:"request_max_time"`
	UpstreamConnectMaxTime  float64 `json:"upstream_connect_max_time"`
	UpstreamHeaderMaxTime   float64 `json:"upstream_header_max_time"`
	UpstreamResponseMaxTime float64 `json:"upstream_response_max_time"`
}

func NewLogAnalyser(myConfig *MyConfig) *LogAnalyser {
	results := []*LogLineAnalysisResult{}
	logHandler := NewLogHandler()
	l := &LogAnalyser{myConfig:myConfig, analysisResults:results, logHandler:logHandler}
	l.setSeek()
	return l
}

func (l *LogAnalyser) Analysis() {
	l.initAnalysisResults()
	l.statistic()
}

func (l *LogAnalyser) setSeek() {
	if !FileExists(l.myConfig.SeekFile) {
		l.seek = 0
	} else {
		file, err := os.OpenFile(l.myConfig.SeekFile, os.O_RDONLY, 0644)
		PanicCheck(err)
		fileReader := bufio.NewReader(file)
		seekString, err := fileReader.ReadString('\n')
		PanicCheck(err)
		l.seek, _ = strconv.ParseInt(strings.Replace(seekString, "\n", "", -1), 10, 64)
	}
}

func (l *LogAnalyser) initAnalysisResults() {

	if !FileExists(l.myConfig.ReadFile) {
		panic(fmt.Sprintf("file not exists %s", l.myConfig.ReadFile))
	}

	openFile, err := os.Open(l.myConfig.ReadFile)
	PanicCheck(err)

	defer openFile.Close()

	logFileReader := bufio.NewReader(openFile)
	PanicCheck(err)

	line := 0

	stats, err := openFile.Stat()
	PanicCheck(err)

	for {
		openFile.Seek(l.seek, 0)
		lineContent, err := logFileReader.ReadString('\n')

		if err == io.EOF {
			break
		}

		WarnCheck(err, "readline error:")

		if lineContent != "" && lineContent != "\n" {
			line++
		}

		l.seek += int64(len(lineContent))

		logLineAnalysisResult, err := l.logHandler.AnalysisLine(lineContent)

		if logLineAnalysisResult != nil {
			l.analysisResults = append(l.analysisResults, logLineAnalysisResult)
		}

		if line >= l.myConfig.ReadLimit {
			break
		}
	}

	if l.myConfig.Mode == "spec" {
		l.seek = stats.Size()
	}
}

func (l *LogAnalyser) statistic() {
	l.saveSeekFile()
	result := &LogAnalysisResult{Count2xx:0, Count3xx:0, Count4xx:0, Count5xx:0, RequestMaxTime:0,
		UpstreamConnectMaxTime:0, UpstreamHeaderAvgTime:0, UpstreamResponseMaxTime:0}
	result.UpdatedTime = time.Now().String()

	var totalRequestTime, totalUpstreamConnectTime, totalUpstreamHeaderTime, totalUpstreamResponseTime float64 = 0, 0,
		0, 0

	for _, logLineAnalysisResult := range l.analysisResults {
		result.addStatus(logLineAnalysisResult.Status)

		totalRequestTime += logLineAnalysisResult.RequestTime
		totalUpstreamConnectTime += logLineAnalysisResult.UpstreamConnectTime
		totalUpstreamHeaderTime += logLineAnalysisResult.UpstreamHeaderTime
		totalUpstreamResponseTime += logLineAnalysisResult.UpstreamResponseTime

		result.RequestMaxTime = result.getMax(result.RequestMaxTime,
			logLineAnalysisResult.RequestTime)
		result.UpstreamConnectMaxTime = result.getMax(result.UpstreamConnectMaxTime,
			logLineAnalysisResult.UpstreamConnectTime)
		result.UpstreamHeaderMaxTime = result.getMax(result.UpstreamHeaderMaxTime,
			logLineAnalysisResult.UpstreamHeaderTime)
		result.UpstreamResponseMaxTime = result.getMax(result.UpstreamResponseMaxTime,
			logLineAnalysisResult.UpstreamResponseTime)
	}

	length := len(l.analysisResults)
	if length != 0 {
		result.RequestAvgTime = totalRequestTime / float64(length)
		result.UpstreamConnectAvgTime = totalUpstreamConnectTime / float64(length)
		result.UpstreamHeaderAvgTime = totalUpstreamHeaderTime / float64(length)
		result.UpstreamResponseAvgTime = totalUpstreamResponseTime / float64(length)
	}
	l.saveResultFile(result)
}

func (l *LogAnalyser) saveSeekFile() {
	file, err := os.OpenFile(l.myConfig.SeekFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	PanicCheck(err)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%d\n", l.seek))
}

func (l *LogAnalyser) saveResultFile(result *LogAnalysisResult) {
	file, err := os.OpenFile(l.myConfig.ResultFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	PanicCheck(err)
	defer file.Close()
	js, _ := json.Marshal(result)
	file.Write(js)
}

func (r *LogAnalysisResult) getMax(time1, time2 float64) float64 {
	return math.Max(time1, time2)
}

func (r *LogAnalysisResult) addStatus(status int) {
	flag := status / 100
	switch flag {
	case 2 :
		r.Count2xx++
	case 3 :
		r.Count3xx++
	case 4 :
		r.Count4xx++
	case 5 :
		r.Count5xx++
	}
}