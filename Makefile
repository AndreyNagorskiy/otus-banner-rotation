BIN := "./bin/banner-rotation"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
GOLANGCI_LINT_VERSION := "v1.63.4"
MIGRATIONS_DIR := "./internal/storage/migrations"
DB_DSN := "postgres://banner_rotation:banner_rotation@localhost:25433/banner_rotation?sslmode=disable"
PROTO_OUT_DIR = pb/api

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/server
run:
	docker compose up

run-bin: build
	$(BIN)

test:
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

lint: install-lint-deps
	golangci-lint run ./...

lint-fix: install-lint-deps
	golangci-lint run ./... --fix

generate-proto:
	rm -rf $(PROTO_OUT_DIR)
	mkdir -p $(PROTO_OUT_DIR)
	protoc --proto_path=api/proto \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		api/proto/bannerrotation/bannerrotation.proto

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) status

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: build run run-bin test lint lint-fix generate-proto migrate-up migrate-down migrate-status migrate-create