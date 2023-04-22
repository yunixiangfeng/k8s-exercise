IMAGE_NAME ?= sidecar-injector

PWD := $(shell pwd)
BASE_DIR := $(shell basename $(PWD))

export GOPATH ?= $(GOPATH_DEFAULT)

IMAGE_TAG ?= $(shell date +v%Y%m%d)-$(shell git describe --match=$(git rev-parse --short-8 HEAD) --tags --always --dirty
build:
	@echo "Building the $(IMAGE_NAME) binary..."
	@CGO_ENABLED=0 go build -o $(IMAGE_NAME) ./pkg/
build-linux:
	@echo "Building the $(IMAGE_NAME) binary for Docker (linux)..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(IMAGE_NAME) ./pkg/
#################################
# image section
#################################

image: build-image

build-image: build-linux
	@echo "Building the docker image: $(IMAGE_NAME)..."
	@docker build -t $(IMAGE_NAME) f Dockerfile .

.PHONY: all build image