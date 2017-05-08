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