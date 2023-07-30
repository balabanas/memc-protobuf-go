# memc-protobuf-go
MemcLoad v2 - Memcached Data Load with GoLang




The protocol buffer compiler requires a plugin to generate Go code. Install it using Go 1.16 or higher by running:


go mod init example.com/mymodule
go get github.com/bradfitz/gomemcache/memcache@v0.0.0-20230611145640-acc696258285
go get google.golang.org/protobuf/cmd/protoc-gen-go@latest



Need to add go package path to .proto file before creation of .go file: 
option go_package = "./";

protoc  --go_out=. ./appsinstalled.proto

For src location outside c:/users/abalabanov/go/src, i've managed to compile in PowerShell only with $env:GOBIN = "c:/Users/abalabanov/go/bin" set first, and then cd c:/projects/memc-protobuf-go, and then:

`protoc  --go_out=./proto ./proto/appsinstalled.proto`


