# memc-protobuf-go
MemcLoad v2 - Memcached Data Load with GoLang

## Goal
Write GoLang utility to load data from .gzipped files into several instances of memcached.

## Run
`go run main.go`


## Setup for dev
```
go mod init github.com/balabanas/memc-protobuf-go
go get github.com/bradfitz/gomemcache/memcache@v0.0.0-20230611145640-acc696258285
go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Need to add go package path to .proto file before creation of .go file: 
`option go_package = "./";`

Create .proto file for Go:
`protoc  --go_out=./proto ./proto/appsinstalled.proto`

## Notes
Mind using PowerShell with `$env:GOBIN = "<path to bin, usually c:/Users/<user>/go/bin>"` before using protoc
