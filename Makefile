WHOAMI ?= $(shell whoami)
CWD := $(shell pwd)
NAME := sneak
BIN_NAME := sneak
INSTALL_LOCATION := /usr/local/bin
COMMIT := $(shell git rev-parse --short HEAD)
TODAY = $(shell date +%Y-%m-%d)
VERSION := $(COMMIT)-$(TODAY)

BUILD_OUTPUT_DIR := $(CWD)/build
BINARY_LOCATION := $(BUILD_OUTPUT_DIR)/$(BIN_NAME)
MODULE := $(shell go list -m)
CMD_MODULE := $(MODULE)/cmd/$(BIN_NAME)
LOCAL_NETWORK := $(shell ifconfig | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $$2}')

${BUILD_OUTPUT_DIR}:
	@mkdir -p $(BUILD_OUTPUT_DIR)

ifneq ($(MAKECMDGOALS),dev,run)
-include env.mk
endif

UNAME_S := $(shell uname -s)
ifeq ($(PLATFORM),)
ifeq ($(UNAME_S),Darwin)
PLATFORM ?= darwin
endif
ifeq ($(UNAME_S),Linux)
PLATFORM ?= linux
endif
endif

GOOS = $(PLATFORM)
GOARCH ?= amd64

GO := $(shell command -v go 2>/dev/null)
GO_LINKER_FLAGS = -X $(CMD_MODULE).Version=$(VERSION)
GO_BUILD_FLAGS = -mod=vendor -a --installsuffix cgo -ldflags "$(GO_LINKER_FLAGS)" -o $(BINARY_LOCATION)

######################################### DEPS

vendor: go.sum tidy
	@GO111MODULE=on $(GO) mod vendor

.PHONY: tidy
tidy:
	@GO111MODULE=on go mod tidy

######################################### BIN

.PHONY: bin
bin: vendor ${BUILD_OUTPUT_DIR} ## build the sneak binary
	@echo "compiling ${NAME}..."
	@export GOOS=$(GOOS) GOARCH=$(GOARCH) && \
		export GO111MODULE=on && \
		export CGO_ENABLED=0 && \
		$(GO) build $(GO_BUILD_FLAGS)
	@echo "${NAME} bin compiled!"

.PHONY: install
install: bin ## install the sneak binary to /usr/local/bin
	@echo "installing sneak to ${INSTALL_LOCATION}"
	@cp ${BINARY_LOCATION} ${INSTALL_LOCATION}
	@chmod 755 ${INSTALL_LOCATION}/${BIN_NAME}
	@echo "installation complete"

######################################### LOCAL DOCKER DEVELOPMENT

.PHONY: dev-bin
dev-bin: vendor ${BUILD_OUTPUT_DIR} ## build the sneak binary for linux for easier dev on running sneaker container locally
	@echo "compiling ${NAME} as linux distro for mounting the bin in local-dev..."
	@export GOOS=linux GOARCH=$(GOARCH) && \
		export GO111MODULE=on && \
		export CGO_ENABLED=0 && \
		$(GO) build $(GO_BUILD_FLAGS)
	@echo "${NAME} bin compiled!"

.PHONY: dev-docker
dev-docker: dev-bin ## build a docker image for local development of sneak binary and sneaker env
	@echo "HACK THE BOC USERNAMEIMPORT $(HTB_USERNAME)"
	@DOCKER_BUILDKIT=1 docker build --build-arg USER=$(HTB_USERNAME) -f Dockerfile.dev -t sneaker .

.PHONY: run
run: dev-docker ## run sneak in a containerized environment
	@docker run \
		--privileged \
		--sysctl net.ipv6.conf.all.disable_ipv6=0 \
		--env LOCAL_NETWORK=$(LOCAL_NETWORK) \
		--cap-add=NET_ADMIN \
		-v $(HOME)/.sneak/:/home/$(HTB_USERNAME)/.sneak \
		-v $(CWD)/build/sneak:/go/bin/sneak \
		-p 8118:8118 \
		-it sneaker \
		 /bin/sh

######################################### DOCKER DIST

.PHONY: image-sneak
image-sneak: ## buids a docker image of only the sneak binary
	@docker build -f dist/sneak/Dockerfile -t docker.io/kfrench/sneak:$(VERSION) .
	@docker tag kfrench/sneak:$(VERSION) docker.io/kfrench/sneak:latest

.PHONY: push-sneak
push-sneak: image-sneak ## pushes the sneak image to docker
	@docker push docker.io/kfrench/sneak:$(VERSION)
	@docker push docker.io/kfrench/sneak:latest

.PHONY: image-sneaker
image-sneaker: ## build the sneaker image
	@DOCKER_BUILDKIT=1 docker build -f dist/sneaker/Dockerfile --ssh default -t docker.io/kfrench/sneaker:$(VERSION) .
	@docker tag docker.io/kfrench/sneaker:$(VERSION) docker.io/kfrench/sneaker:latest

.PHONY: push-sneaker
push-sneaker: image-sneaker ## push the sneaker image
	@docker push docker.io/kfrench/sneaker:$(VERSION)
	@docker push docker.io/kfrench/sneaker:latest

######################################### RUN IT

.PHONY: pull
pull: ## pull the latest sneaker image
	@docker pull docker.io/kfrench/sneaker:latest

.PHONY: up
up: pull ## run sneak's containerized environment, sneaker, locally
	@docker run \
		--privileged \
		--sysctl net.ipv6.conf.all.disable_ipv6=0 \
		--env LOCAL_NETWORK=$(LOCAL_NETWORK) \
		--cap-add=NET_ADMIN \
		-p 8118:8118 \
		-it docker.io/kfrench/sneaker:latest \
		 /bin/sh

######################################### TESTING

.PHONY: lint
lint: ## go linter and shadow tool
	@$(GO) get -u golang.org/x/lint/golint
	@$(GO) get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	@$(GO) vet $(shell go list ./... | grep -v 'vendor')
	@$(GO) vet -vettool=$(shell which shadow) ./...

.PHONY: test
test: lint ## run linter and unit tests 
	@echo "running tests..."
	@$(GO) test ./...

.PHONY: dbweb
dbweb: ## view database info on localhost:8080
	@go get -u github.com/evnix/boltdbweb
	@boltdbweb -d $(HOME)/.sneak/sneak.db

######################################### UTILS

.PHONY: clean
clean: ## delete the build binary
	@rm -rf ${BUILD_OUTPUT_DIR}
	@echo "removed ${BUILD_OUTPUT_DIR}..."

.PHONY: help
help: ## lists some available makefile commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
