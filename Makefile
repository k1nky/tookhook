SHELL:=/bin/bash
STATICCHECK=$(shell which staticcheck)
PLUGINS:=$(wildcard plugins/*)
BUILD_PATH:=build
MKFILE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR:=$(dir $(MKFILE_PATH))
BULID_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date +'%Y/%m/%d %H:%M:%S')
BUILD_VERSION=$(shell git tag --points-at HEAD)
LDFLAGS=-ldflags "-X main.buildVersion=${BUILD_VERSION} -X main.buildCommit=${BULID_COMMIT} -X 'main.buildDate=${BUILD_DATE}'"


.DEFAULT_GOAL := build

define pluginsmake
	for d in ${PLUGINS}; do \
		make -C $$d $(1) BUILD_PATH=${MKFILE_DIR}${BUILD_PATH}; \
	done
endef

test:
	go test -cover ./...

vet:
	go vet ./...
	$(STATICCHECK) ./...

generate:
	go generate ./...
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/plugin/proto/*.proto

gvt: generate vet test

cover:
	go test -cover ./... -coverprofile cover.out
	go tool cover -html cover.out -o cover.html

clean:
	rm -rf build
	rm -rf dev
	rm -rf plugins
	rm -f cover.*

build: gvt 
	CGO_ENABLED=0 go build -o ${BUILD_PATH}/tookhook ${LDFLAGS} cmd/*.go

plugin:
	$(call pluginsmake,build)

plugin-dev:
	$(call pluginsmake,plugin-dev)

docker:
	docker build -t k1nky/tookhook:latest .

run:
	go run ./cmd

addplugins:
	git submodule add --force --name telegram git@github.com:k1nky/tookhook-plugin-telegram.git plugins/telegram
	git submodule add --force --name pachca git@github.com:k1nky/tookhook-plugin-pachca.git plugins/pachca
	git submodule update --init --recursive --remote

prepare:
	go mod tidy	
	go install go.uber.org/mock/mockgen@v0.4.0
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.7
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@v0.7.7