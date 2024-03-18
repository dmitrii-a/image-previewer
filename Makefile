BIN_PREVIEWER ?= "./bin/previewer"
DOCKER_IMG ?= "previewer:develop"
GIT_HASH ?= $(shell git log --format="%h" -n 1)
LDFLAGS ?= -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up:
	cd deployments && docker-compose up -d
docker-build:
	cd deployments && docker-compose build
down:
	cd deployments && docker-compose down
restart:
	cd deployments && docker-compose restart
rm:
	cd deployments && docker-compose rm -v

build-image-previewer:
	go build -v -o $(BIN_PREVIEWER) -ldflags "$(LDFLAGS)"  ./cmd/app

build: build-image-previewer

run-image-previewer: build-image-previewer
	$(BIN_PREVIEWER) -config ./configs/config.yaml

run: run-image-previewer

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_PREVIEWER) version

test:
	go test -race ./internal/... ./pkg/...

integration-tests: down
	cd deployments && COMPOSE_PROJECT_NAME=test docker-compose up -d --build
	trap 'cd deployments && COMPOSE_PROJECT_NAME=test docker-compose down -v' EXIT; \
	go test -v ./tests/integration/...

integration-tests-local:
	go test -v ./tests/integration/...

install-lint-deps:
	(which golangci-lint > /dev/null) || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

install-wsl-deps:
	(which wsl > /dev/null) || go install github.com/bombsimon/wsl/v4/cmd/wsl@v4.2.1

lint: install-lint-deps
	golangci-lint run ./...

fmt: install-lint-deps install-wsl-deps
	go fmt ./...
	gofumpt -w -l -extra .
	golines -w .
	wsl --fix ./...
	golangci-lint run ./... --fix

install-mockery:
	go install github.com/vektra/mockery/v2@v2.40.1

generate-mocks:
	mockery --output=./tests/mocks --exclude=vendor --all
