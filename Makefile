.PHONY: build

C2J_VERSION=v0.0.1
C2J_GIT_COMMIT=$(shell git rev-parse --short HEAD)
C2J_DATE=$(shell date)
C2J_PACKAGE="go.krishnaiyer.dev/csvtojson"

test:
	go test ./... -cover

build.local:
	go build \
	-ldflags="-X '${C2J_PACKAGE}/cmd.version=${C2J_VERSION}' \
	-X '${C2J_PACKAGE}/cmd.gitCommit=${C2J_GIT_COMMIT}' \
	-X '${C2J_PACKAGE}/cmd.buildDate=${C2J_DATE}'" main.go

build.dist:
	goreleaser --snapshot --skip-publish --rm-dist
clean:
	rm -rf dist
