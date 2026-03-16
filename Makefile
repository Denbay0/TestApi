APP=edge-api

.PHONY: proto run build test docker

proto:
	PATH="$(shell go env GOPATH)/bin:$(PATH)" buf generate

run:
	go run ./cmd/edge-api

build:
	CGO_ENABLED=0 go build -o bin/$(APP) ./cmd/edge-api

test:
	go test ./...

docker:
	docker build -t $(APP):local .
