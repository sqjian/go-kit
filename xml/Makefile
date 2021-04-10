IMG:='registry.cn-hangzhou.aliyuncs.com/sqjian/venv:ubuntu20_04'

.PHONY: build env

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

build: dep
	go build -v
venv:
	docker run \
			-it \
			--rm \
			--net=host \
			--name=venv \
			-v ${PWD}:/lab \
			-w /lab \
			${IMG} bash
dep:
	go mod tidy
	apt-get install -y pkg-config
	apt-get install -y libxml2-dev
