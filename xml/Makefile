IMG:='registry.cn-hangzhou.aliyuncs.com/sqjian/venv:ubuntu20_04'

.PHONY: build env

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

build: env
  go build -v
env:
  go mod tidy
  apt-get install -y pkg-config
  apt-get install -y libxml2-dev