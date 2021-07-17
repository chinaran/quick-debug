.PHONY: all build

all: proto-gen build

proto-gen:
	protoc \
		-I. \
		-I="${GOPATH}/src" \
		-I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
		--gogofaster_out=plugins=grpc:. \
		./proto/debug/*.proto \

build:
	go vet ./cmd/quick-debug
	CGO_ENABLED=0 GOOS=linux go build -o quick-debug ./cmd/quick-debug
	upx quick-debug

local:
	go vet ./cmd/quick-debug
	go build -o quick-debug ./cmd/quick-debug
	go vet ./cmd/quick-debug-client
	go build -o quick-debug-client ./cmd/quick-debug-client