MIGRATIONS_DIR := "./internal/storage/migrations"
DB_DSN := "postgres://user:password@localhost:5432/banner_rotation?sslmode=disable"

OUT_DIR = pb/api

test:
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

lint-fix: install-lint-deps
	golangci-lint run ./... --fix

generate-proto:
	rm -rf $(OUT_DIR)
	mkdir -p $(OUT_DIR)
	protoc --proto_path=api/proto \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		api/proto/bannerrotation/bannerrotation.proto

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) status

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: test lint lint-fix generate-proto migrate-up migrate-down migrate-status migrate-create