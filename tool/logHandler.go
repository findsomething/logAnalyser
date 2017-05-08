package tool

import (
	"github.com/satyrius/gonx"
	"strings"
	"strconv"
)

type LogHandler struct {

}

type LogLineAnalysisResult struct {
	RemoteAddr           string
	RemoteUser           string
	TimeLocal            string
	Request              string
	Status               int
	BodyBytesSent        string
	HttpRefer            string
	HttpUserAgent        string
	RequestTime          float64
	UpstreamConnectTime  float64
	UpstreamHeaderTime   float64
	UpstreamResponseTime float64
}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (l *LogHandler) AnalysisLine(line string) (*LogLineAnalysisResult, error) {
	logReader := strings.NewReader(line)
	nginxConfig := strings.NewReader(`
		log_format main '$remote_addr - $remote_user [$time_local] '
                             '"$request" $status $body_bytes_sent '
                             '"$http_referer" "$http_user_agent" '
                             'rt="$request_time" uct="$upstream_connect_time" uht="$upstream_header_time" urt="$upstream_response_time"'
	`)
	reader, _ := gonx.NewNginxReader(logReader, nginxConfig, "main")
	rec, err := reader.Read()
	if err != nil {
		return nil, err
	}

	result := l.getResultFromEntry(rec)

	return result, nil
}

func (l *LogHandler) getResultFromEntry(entry *gonx.Entry) *LogLineAnalysisResult {
	r := &LogLineAnalysisResult{}

	statusString, _ := entry.Field("status")

	status, _ := strconv.Atoi(statusString)

	r.RemoteAddr, _ = entry.Field("remote_addr")
	r.RemoteUser, _ = entry.Field("remote_user")
	r.TimeLocal, _ = entry.Field("time_local")
	r.Request, _ = entry.Field("request")
	r.Status = status
	r.BodyBytesSent, _ = entry.Field("body_bytes_sent")
	r.HttpRefer, _ = entry.Field("http_referer")
	r.HttpUserAgent, _ = entry.Field("http_user_agent")
	r.RequestTime, _ = entry.FloatField("request_time")
	r.UpstreamConnectTime, _ = entry.FloatField("upstream_connect_time")
	r.UpstreamHeaderTime, _ = entry.FloatField("upstream_header_time")
	r.UpstreamResponseTime, _ = entry.FloatField("upstream_response_time")

	return r
}

