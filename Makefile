.PHONY: all build

build:
	CGO_ENABLED=0 GOOS=linux go build -o quick-debug *.go
	upx quick-debug