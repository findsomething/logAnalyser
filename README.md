# logAnalyser
nginx log analyser

## install 
```
git clone git@github.com:findsomething/logAnalyser.git

go get github.com/droundy/goopt
go get github.com/satyrius/gonx
```

## usage
```
go build main.go

./main --input=thePathToNginxAccessLogFile --output=outputFileName 
```