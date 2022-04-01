.PHONY: all build local lint

all: proto-gen build

proto-gen:
	protoc \
		-I. \
		-I="${GOPATH}/src" \
		-I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
		--gogofaster_out=plugins=grpc:. \
		./proto/debug/*.proto \

lint:
	@# golangci-lint run ./...
	revive -formatter friendly ./...

build:
	go vet ./cmd/quick-debug
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/quick-debug ./cmd/quick-debug
	go vet ./cmd/quick-debug-client
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/quick-debug-client ./cmd/quick-debug-client

install: local

local:
	go vet ./cmd/quick-debug
	go install ./cmd/quick-debug
	go vet ./cmd/quick-debug-client
	go install ./cmd/quick-debug-client
