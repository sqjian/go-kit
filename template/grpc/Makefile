.PHONY: idl env

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

idl:
	cd idl && make
env:
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/golang/protobuf/proto
	go get google.golang.org/grpc
	go get github.com/envoyproxy/protoc-gen-validate
	go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/gogo/protobuf/protoc-gen-gogofast
	go get github.com/gogo/protobuf/protoc-gen-gogofaster
	go get github.com/gogo/protobuf/protoc-gen-gogoslick
	go get github.com/gogo/protobuf/gogoproto