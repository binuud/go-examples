# Maintainer Binu Udayakumar <binu@dronasys.com>

REPO=binuud


CONSCIOUS_APP_NAME= go-examples

DEBUG_IMAGE= debug-golang

BUILD_VER= a1.0.0
DOCKER_HUB_TAG= a1.0.0

GOLANG_IMAGE=golang:1.27.2

TAGGED_NAME = $(REPO)/$(APP_NAME)

MULTI_PLATFORM_DOCKER = --platform=linux/amd64,linux/arm64

# Optional environmental variables

BASE_DIR	:=	$(shell git rev-parse --show-toplevel)
BIN			:=	$(BASE_DIR)/bin

## Needs protoc to be installed

ifndef PROTOC
PROTOC = protoc
endif

ifndef PROTOWEB
PROTOWEB = protoc-gen-grpc-web
endif

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help vendor vendor-update
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

##
## Build, Vendor, Publish
##

main: ## Build main
	go build $(BUILDFLAGS) -o main .


tidy: ## run tidy
	go mod tidy

vendor: ## run vendor
	go mod tidy
	go mod vendor

staticcheck:
	GOFLAGS="" go install honnef.co/go/tools/cmd/staticcheck@v0.2.2
	staticcheck $(PKGS)

errcheck:
	go get -u github.com/kisielk/errcheck
	errcheck -ignoregenerated -verbose -blank $(PKGS)

##
## Test Code
##

test: ## Run test 
	./scripts/runTest.sh

start-dev-container: ## start golang container
	docker run --rm -it --name $(CONSCIOUS_APP_NAME) -p 9080:9080 -p 9090:9090 -v .:/app -w /app dronasys-com/golang:1.27.1

connect-dev-container: ## connect to dev container
	docker exec -it $(CONSCIOUS_APP_NAME) bash
