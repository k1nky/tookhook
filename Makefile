SHELL:=/bin/bash
STATICCHECK=$(shell which staticcheck)
PLUGINS:=$(wildcard plugins/*)
BUILD_PATH:=build
MKFILE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR:=$(dir $(MKFILE_PATH))


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
	CGO_ENABLED=0 go build -o ${BUILD_PATH}/tookhook cmd/*.go

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
	go install go.uber.org/mock/mockgen@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest