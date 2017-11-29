export GOPATH=$GOPATH:`pwd`
go get
go build -o bin/tenchmark tenchmark/*.go
go build -o bin/ting ting/*.go
go build -o bin/toxy toxy.go
