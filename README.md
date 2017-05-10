# logAnalyser
nginx log analyser

## install 
```
git clone git@github.com:findsomething/logAnalyser.git

go get github.com/droundy/goopt
go get github.com/satyrius/gonx
```

## notice
log format 
```
log_format main '$remote_addr - $remote_user [$time_local] '
                             '"$request" $status $body_bytes_sent '
                             '"$http_referer" "$http_user_agent" '
                             'rt="$request_time" uct="$upstream_connect_time" uht="$upstream_header_time" urt="$upstream_response_time"';
```

## usage
```
go build main.go

./main --input=thePathToNginxAccessLogFile --outputName=outputFileName --outputPath=/var/tmp/output --mode=spec
```

## result

outputFileName content:
```
{
    "update_time":1494385717,
    "1xx_count":1,
    "2xx_count":8,
    "3xx_count":0,
    "4xx_count":0,
    "5xx_count":2,
    "request_avg_time":0.5931,
    "upstream_connect_avg_time":0.0007,
    "upstream_header_avg_time":0.5931,
    "upstream_response_avg_time":0.5931,
    "request_max_time":1.697,
    "upstream_connect_max_time":0.002,
    "upstream_header_max_time":1.697,
    "upstream_response_max_time":1.697
}
```