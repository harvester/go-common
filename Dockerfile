FROM registry.suse.com/bci/golang:1.25 AS builder
ENV GOTOOLCHAIN=auto

RUN zypper -n rm container-suseconnect && \
    zypper -n install git curl gzip tar wget awk && \
    zypper -n clean -a

## install golangci
COPY --from=golangci/golangci-lint:v2.11.4-alpine@sha256:72bcd68512b4e27540dd3a778a1b7afd45759d8145cfb3c089f1d7af53e718e9 \
    /usr/bin/golangci-lint /usr/local/bin/golangci-lint

WORKDIR /go/src/github.com/harvester/go-common/
