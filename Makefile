SHELL:=/bin/bash
STATICCHECK=$(shell which staticcheck)

.DEFAULT_GOAL := build

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

build: gvt 
	CGO_ENABLED=0 go build -o tookhook cmd/*.go

plugin:
	go build -o pachca plugins/pachca/*.go

run:
	go run ./cmd

prepare:
	go mod tidy
	go install go.uber.org/mock/mockgen@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest