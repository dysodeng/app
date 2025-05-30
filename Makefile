SHELL:=/bin/bash

DOCKER_REGISTRY_NAME=registry.cn-chengdu.aliyuncs.com/dysodeng/app
PLATFORM_ARCH=linux/amd64
PROD_DOCKER_REGISTRY_NAME=registry.cn-chengdu.aliyuncs.com/dysodeng/app-release
PROD_PLATFORM_ARCH=linux/amd64
IMAGE_VERSION=$(shell date +"%Y%m%d%H%M")

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

build:
	docker buildx build --platform ${PLATFORM_ARCH} -t ${DOCKER_REGISTRY_NAME}:${IMAGE_VERSION} .
	docker push ${DOCKER_REGISTRY_NAME}:${IMAGE_VERSION}

build-prod:
	docker buildx build --platform ${PROD_PLATFORM_ARCH} -t ${PROD_DOCKER_REGISTRY_NAME}:${IMAGE_VERSION} .
	docker push ${PROD_DOCKER_REGISTRY_NAME}:${IMAGE_VERSION}

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
