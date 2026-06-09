ROOT := $(realpath $(dir $(realpath $(firstword $(MAKEFILE_LIST)))))

export DOCKER_BUILDKIT := 1

BUILD_IMAGE ?= harvester-go-common-builder

ifdef CI
  BOLD  :=
  CYAN  :=
  RESET :=
else
  BOLD  := \033[1m
  CYAN  := \033[36m
  RESET := \033[0m
endif

BANNER = @printf "$(BOLD)$(CYAN)[target: $@]$(RESET)\n"

MK_DOCKER_PROGRESS ?= plain

DOCKER_BUILD = docker build \
	--progress=$(MK_DOCKER_PROGRESS) \
	--target builder \
	-t $(BUILD_IMAGE) \
	-f $(ROOT)/Dockerfile $(ROOT)

DOCKER_RUN = docker run --rm \
	--user $(shell id -u):$(shell id -g) \
	-e HOME=/tmp \
	-e GOPATH=/tmp/go \
	-e GOCACHE=/tmp/go/cache \
	-e GOMODCACHE=/tmp/go/pkg/mod \
	-w /go/src/github.com/harvester/go-common/ \
	-v $(ROOT):/go/src/github.com/harvester/go-common/ \
	$(BUILD_IMAGE)

.PHONY: build ci clean-all default test validate

default: ci

ci: validate test

build:
	$(BANNER)
	$(DOCKER_BUILD)
	$(DOCKER_RUN) go build -v ./...

test:
	$(BANNER)
	$(DOCKER_BUILD)
	$(DOCKER_RUN) go test -cover -tags=test ./...

validate:
	$(BANNER)
	$(DOCKER_BUILD)
	$(DOCKER_RUN) golangci-lint run
	$(DOCKER_RUN) sh -c "test -z \"\$$(gofmt -l . | tee /dev/stderr)\""

clean-all:
	$(BANNER)
	docker rmi $(BUILD_IMAGE) || true

.DEFAULT_GOAL := default
