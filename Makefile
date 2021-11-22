USERNAME := alekseevdenis
APP_NAME := shortener
VERSION := latest
PROJECT := github.com/alekceev/go-shortener
GIT_COMMIT := $(shell git rev-parse HEAD)

check:
	golangci-lint run -c golangci-lint.yaml

test:
	go test ./...

generate:
	go generate ./...

run:
	go install -ldflags="-X '$(PROJECT)/app/config.Version=$(VERSION)' \
	-X '$(PROJECT)/app/config.Commit=$(GIT_COMMIT)'" ./cmd/shorturl && shorturl

build_container:
	docker build --build-arg=GIT_COMMIT=$(GIT_COMMIT) --build-arg=VERSION=$(VERSION) --build-arg=PROJECT=$(PROJECT) \
	-t docker.io/$(USERNAME)/$(APP_NAME):$(VERSION) .
