# Builder
FROM golang:1.16 AS builder

ADD . /go/src
WORKDIR /go/src
RUN GOPROXY=https://goproxy.io go mod download
RUN make build
RUN apt update && apt install -y upx
RUN upx /go/src/quick-debug
RUN upx /go/src/quick-debug-client

# Runner
FROM alpine:3.13.5

RUN apk --no-cache add curl tzdata

ENV TZ Asia/Shanghai

COPY --from=builder /go/src/quick-debug /usr/local/bin/
COPY --from=builder /go/src/quick-debug-client /usr/local/bin/
