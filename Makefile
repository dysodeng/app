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
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
	fi
	@if [ ! -L proto ]; then \
		ln -sf api/proto proto; \
		echo "Created proto symlink for IDE support"; \
	fi
	@if [ -d .idea ] && [ ! -f .idea/protobuf.xml ]; then \
    		cp ./scripts/setting/protobuf.xml ./.idea/protobuf.xml; \
    		echo "Created .idea/protobuf.xml for GoLand support"; \
    fi
	@go mod tidy
	@go mod download
	@echo "Initialized successfully!"

# Protobuf相关命令
.PHONY: proto-gen proto-lint proto-format proto-breaking proto-clean

# 生成protobuf代码
proto-gen:
	@echo "Generating protobuf code..."
	@cd api && buf generate

# 检查protobuf代码规范
proto-lint:
	@echo "Linting protobuf files..."
	@cd api && buf lint

# 格式化protobuf文件
proto-format:
	@echo "Formatting protobuf files..."
	@cd api && buf format -w

# 检查breaking changes
proto-breaking:
	@echo "Checking for breaking changes..."
	@cd api && buf breaking --against '.git#branch=main'

# 清理生成的代码
proto-clean:
	@echo "Cleaning generated protobuf code..."
	@rm -rf api/generated/go/*

# 清理IDE支持文件
.PHONY: clean-ide
clean-ide:
	@echo "Cleaning IDE support files..."
	@if [ -L proto ]; then rm proto; echo "Removed proto symlink"; fi
	@if [ -f .idea/protobuf.xml ]; then rm .idea/protobuf.xml; echo "Removed .idea/protobuf.xml"; fi

# 完整清理
.PHONY: clean
clean: proto-clean clean-ide
	@echo "Cleaned all generated files and IDE support files"

# 完整的protobuf工作流
proto: proto-lint proto-format proto-gen
	@echo "Protobuf workflow completed!"
