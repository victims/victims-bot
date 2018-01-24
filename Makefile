# Passed into main.go at build time
VERSION := $(shell cat ./VERSION)
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date +%s)

# Used in tagging images
IMAGE_VERSION_TAG := victims-bot:$(VERSION)
IMAGE_DATE_TAG := victims-bot:$(BUILD_TIME)

# Used during all builds
LDFLAGS := -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildTime=${BUILD_TIME}

.PHONY: help clean victims-bot image

default: help

help:
	@echo "Targets:"
	@echo " deps: Install dependencies with govendor"
	@echo "	victims-bot: Builds a victims-bot binary"
	@echo "	clean: cleans up and removes built files"
	@echo "	image: builds a container image"

deps:
	go get github.com/kardianos/govendor
	govendor sync

victims-bot:
	govendor build -ldflags '${LDFLAGS}' -o victims-bot main.go

static-victims-bot:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 govendor build --ldflags '-extldflags "-static" ${LDFLAGS}' -a -o victims-bot main.go

clean:
	go clean
	rm -f victims-bot

image: clean deps static-victims-bot
	sudo docker build -t $(IMAGE_VERSION_TAG) -t $(IMAGE_DATE_TAG) .

autobuild: clean
	sudo docker build -t release -t $(IMAGE_DATE_TAG) -f Dockerfile.autobuild .

test: clean deps
	govendor test -v -cover github.com/victims/victims-bot/process github.com/victims/victims-bot/cmd github.com/victims/victims-bot/web

gofmt:
	gofmt -l cmd/ log/ web/ main.go

golint:
	go get github.com/golang/lint/golint
	golint cmd/ log/ web/

lint: gofmt golint
