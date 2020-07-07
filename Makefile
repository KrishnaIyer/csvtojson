.PHONY: build

VERSION=v0.0.1
GIT_COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date)
PACKAGE="go.krishnaiyer.dev/csvtojson"

test:
	go test ./... -cover

build:
	go build \
	-ldflags="-X '${PACKAGE}/cmd.version=${VERSION}' \
	-X '${PACKAGE}/cmd.gitCommit=${GIT_COMMIT}' \
	-X '${PACKAGE}/cmd.buildDate=${DATE}'" main.go
