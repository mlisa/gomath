# gomath

## Dependencies 

`protoactor-go`
`protobuf`
`protoc-gen-gogoslick`


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
