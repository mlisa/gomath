# gomath

## Installation

`go get github.com/mlisa/gomath/...`

Downloads and compiles both peer and coordinator package (binary in $GOPATH/bin/).

To compile a module run `go run *.go` inside the folder.

## Dependencies 

`protoactor-go`

`protobuf`

`protoc-gen-gogoslick`

`gocui`

`go-cache`

`pigeon`

## Proto.Actor

### Link

https://github.com/AsynkronIT/protoactor-go

### Installation

```
go get github.com/AsynkronIT/protoactor-go/...
cd $GOPATH/src/github.com/AsynkronIT/protoactor-go
go get ./...
make
```

## Protobuf

### Link

https://github.com/google/protobuf

### Installation

Installed by system package manager


## Protoc-gen-gogoslick

### Link

https://github.com/gogo/protobuf

### Installation

```
go get github.com/gogo/protobuf/proto                                                                                                                                     
go get github.com/gogo/protobuf/protoc-gen-gogoslick
go get github.com/gogo/protobuf/gogoproto
```

### Code generation

`protoc -I=. -I=$GOPATH/src --gogoslick_out=. *.proto`

## Pigeon

### Link

https://github.com/mna/pigeon

### Installation

`go get -u github.com/mna/pigeon`

### Code generation

`pigeon -optimize-grammar pegmatch.peg > pegmatch.go`

## gocui

### Link

https://github.com/jroimartin/gocui

### Installation

`go get github.com/jroimartin/gocui`

## go-cache

### Link

https://github.com/patrickmn/go-cache

### Installation

`go get github.com/patrickmn/go-cache`
