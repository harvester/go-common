FROM registry.suse.com/bci/golang:1.25.7 AS builder

RUN zypper -n rm container-suseconnect && \
    zypper -n install git curl gzip tar wget awk && \
    zypper -n clean -a

## install golangci
COPY --from=golangci/golangci-lint:v2.12.2-alpine@sha256:91b27804074a0bacea298707f016911e60cf0cdbc6c7bf5ccacb5f0606d18d60 \
    /usr/bin/golangci-lint /usr/local/bin/golangci-lint

WORKDIR /go/src/github.com/harvester/go-common/
