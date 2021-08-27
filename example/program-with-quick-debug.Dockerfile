# Your Builder
FROM golang:1.16 AS builder

# do some build



# Quick Qebug Tool
FROM ghcr.io/chinaran/quick-debug:0.2-alpine3.13 AS debugger

# Your Runner
FROM alpine:3.13.5

# add quick-debug
COPY --from=debugger /usr/local/bin/quick-debug /usr/local/bin/

# something else
