SHELL := /bin/bash
GIT_VERSION := $(shell git describe --dirty --always --tags)

BUILDTIME=$(shell date +%s)
SCRIPTPATH=$(shell pwd -P)
DRONE_TAG?=$(GIT_VERSION)

help:	## Lists all available commands and a brief description.
	@grep -E '^[a-zA-Z/_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := default

.PHONY: vendor
default: vendor protos/ender.pb.go ## Builds ender
	go build -ldflags="-X 'github.com/kaidyth/ender/command.version=\"$(GIT_VERSION)\"' -X 'github.com/kaidyth/ender/command.architecture=\"$(shell uname)/$(shell arch)\"'" \

protos/ender.pb.go:
	PATH=$PATH:./vendor/bin GOBIN=./vendor/bin ./vendor/bin/protoc --proto_path=./protobuf --go_out=./protos --go_opt=paths=source_relative --go-grpc_out=./protos --go-grpc_opt=paths=source_relative ender.proto

vendor/bin/protoc-gen-go-grpc:
	GOBIN=$(shell pwd)/vendor/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

vendor/bin/protoc:
	curl -Lqss https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-linux-aarch_64.zip -o ./vendor/bin/protoc.zip
	unzip ./vendor/bin/protoc.zip -d ./vendor/bin
	cp ./vendor/bin/bin/protoc ./vendor/bin

vendor/bin/protoc-gen-go: # Installs the protoc binary
	GOBIN=$(shell pwd)/vendor/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

vendor:  ## Go Vendor
	go mod vendor
	go mod tidy
	make  vendor/bin/protoc-gen-go vendor/bin/protoc vendor/bin/protoc-gen-go-grpc
