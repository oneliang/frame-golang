GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)


.PHONY: build
# build
build:
	mkdir -p bin/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
