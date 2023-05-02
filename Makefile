GOPATH:=$(shell go env GOPATH)
export DOCKER_BUILDKIT = 1
export COMPOSE_DOCKER_CLI_BUILD = 1
# include .makeenv

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

.PHONY: proto
proto:
	@protoc --proto_path=. --micro_out=. --go_out=:. proto/connector-callback-service.proto

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o connector-callback-service *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t connector-callback-service:latest --no-cache --build-arg GIT_ENERGY_USERNAME=Mamadues --build-arg GIT_ENERGY_PASSWORD=P3N-c-QxoxxFpBp7n2PJ . 
