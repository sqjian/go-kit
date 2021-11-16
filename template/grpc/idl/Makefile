.PHONY: all go

all:go


go:
	protoc \
      -I . \
      --go_out=plugins=grpc:../ \
      *.proto

