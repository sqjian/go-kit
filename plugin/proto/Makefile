.PHONY: all env

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

all:
	protoc \
		-I . \
		--go_out=":./../" \
		msg.proto