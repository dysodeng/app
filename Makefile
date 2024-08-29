SHELL:=/bin/bash

DOCKER_REGISTRY_NAME=registry.huaxisy.com/ai-research/smart-ward/api
PLATFORM_ARCH=linux/arm64
PROD_DOCKER_REGISTRY_NAME=mcpregistry.huaxishuyi.com/ai/ai-assistant-api
PROD_PLATFORM_ARCH=linux/amd64
IMAGE_VERSION=$(shell date +"%Y%m%d%H%M")

.PHONY: run lint build

# go install github.com/pilu/fresh 全局安装fresh命令,热更新代码
run:
	cp -f ./.env ./var/tmp/.env && fresh -c dev-run.conf

lint:
	go fmt ./internal/... && go fmt ./server/... && go fmt ./main.go

build:
	docker buildx build --platform ${PLATFORM_ARCH} -t ${DOCKER_REGISTRY_NAME}:${IMAGE_VERSION} .
	docker push ${DOCKER_REGISTRY_NAME}:${IMAGE_VERSION}

build-prod:
	docker buildx build --platform ${PROD_PLATFORM_ARCH} -t ${PROD_DOCKER_REGISTRY_NAME}:${IMAGE_VERSION} .
	docker push ${PROD_DOCKER_REGISTRY_NAME}:${IMAGE_VERSION}
