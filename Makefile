SHELL:=/bin/bash

.PHONY: run

# go get github.com/pilu/fresh 全局安装fresh命令,热更新代码
run:wire
	cp -f ./.env ./var/tmp/.env && fresh -c dev-run.conf

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

wire:
	wire ./internal/di

wire-show:
	wire show ./internal/di

install-hooks:
	@chmod +x ./scripts/install-hooks.sh
	@./scripts/install-hooks.sh

.PHONY: init
init: install-hooks
	@echo "Initializing..."
	@go install github.com/pilu/fresh@latest
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
	fi
	@go mod tidy
	@go mod download
	@echo "Initialized successfully!"
