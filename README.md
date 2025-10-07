# Go App Service 项目模板

一个基于 Go 语言的现代化微服务项目模板，采用领域驱动设计（DDD）架构，集成了常用的中间件和工具，可快速启动新的服务项目。

## 🚀 特性

- **现代化架构**: 采用领域驱动设计（DDD）和清洁架构
- **多协议支持**: HTTP、gRPC、WebSocket 多协议服务
- **依赖注入**: 使用 Google Wire 进行依赖注入
- **数据库支持**: 支持 PostgreSQL、MySQL，集成 GORM ORM
- **缓存支持**: Redis 缓存，支持单机、集群、哨兵模式
- **消息队列**: 支持 AMQP 消息队列
- **可观测性**: 集成 OpenTelemetry，支持链路追踪、指标监控、日志收集
- **文件存储**: 支持本地存储、MinIO 对象存储
- **第三方集成**: 微信小程序、阿里云服务等
- **开发工具**: 热重载、代码检查、Git hooks

## 📁 项目结构

```
├── cmd/                    # 应用程序入口
│   └── app/               # 主应用
├── configs/               # 配置文件
├── internal/              # 内部代码
│   ├── application/       # 应用层
│   │   ├── file/         # 文件服务
│   │   └── passport/     # 认证服务
│   ├── di/               # 依赖注入
│   ├── domain/           # 领域层
│   │   ├── file/        # 文件领域
│   │   ├── passport/    # 认证领域
│   │   ├── shared/      # 共享领域
│   │   └── user/        # 用户领域
│   ├── infrastructure/   # 基础设施层
│   │   ├── config/      # 配置管理
│   │   ├── event/       # 事件处理
│   │   ├── migration/   # 数据库迁移
│   │   ├── persistence/ # 数据持久化
│   │   ├── server/      # 服务器
│   │   └── shared/      # 共享基础设施
│   └── interfaces/       # 接口层
│       ├── grpc/        # gRPC 接口
│       ├── http/        # HTTP 接口
│       └── websocket/   # WebSocket 接口
├── scripts/              # 脚本文件
├── uploads/              # 上传文件目录
└── var/                  # 运行时文件
    ├── data/            # 数据文件
    ├── logs/            # 日志文件
    └── tmp/             # 临时文件
```

## 🛠️ 技术栈

### 核心框架
- **Go**: 1.24+
- **Gin**: HTTP Web 框架
- **gRPC**: 高性能 RPC 框架
- **WebSocket**: 实时通信

### 数据存储
- **GORM**: ORM 框架
- **PostgreSQL/MySQL**: 关系型数据库
- **Redis**: 缓存和会话存储

### 中间件和工具
- **Google Wire**: 依赖注入
- **Viper**: 配置管理
- **Zap**: 结构化日志
- **JWT**: 身份认证
- **OpenTelemetry**: 可观测性
- **AMQP**: 消息队列

### 第三方服务
- **微信小程序**: 微信生态集成
- **阿里云**: 云服务集成
- **MinIO**: 对象存储

## 🚀 快速开始

### 环境要求

- Go 1.24+
- PostgreSQL 或 MySQL
- Redis
- RabbitMQ (可选)

### 安装

1. **克隆项目**
```bash
git clone <your-repo-url>
cd app-service
```

2. **初始化项目**
```bash
make init
```
这个命令会：
- 安装 Git hooks
- 安装 fresh 热重载工具
- 复制 `.env.example` 到 `.env`
- 下载依赖

3. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库、Redis 等连接信息
```

4. **生成依赖注入代码**
```bash
make wire
```

### 运行

#### 开发模式（热重载）
```bash
make run
```

#### 生产模式
```bash
go build -o app main.go
./app
```

### 测试

```bash
# 运行测试
make test

# 代码检查
make lint
```

## ⚙️ 配置说明

### 环境变量配置

主要的环境变量配置项：

```bash
# 应用配置
APP_ENV=dev                    # 环境：dev/prod
APP_DEBUG=true                 # 调试模式
APP_DOMAIN=http://127.0.0.1:8080

# 服务端口
SERVER_HTTP_PORT=8080          # HTTP 服务端口
SERVER_GRPC_PORT=3000          # gRPC 服务端口
SERVER_HEALTH_CHECK_PORT=5000  # 健康检查端口

# 数据库配置
MAIN_DB_HOST=localhost
MAIN_DB_PORT=5432
MAIN_DB_DATABASE=app
MAIN_DB_USERNAME=postgres
MAIN_DB_PASSWORD=your_password

# Redis 配置
MAIN_REDIS_HOST=127.0.0.1
MAIN_REDIS_PORT=6379
MAIN_REDIS_PASSWORD=
MAIN_REDIS_KEY_PREFIX=app

# JWT 配置
SECURITY_JWT_SECRET=your-secret-key

# 可观测性配置
MONITOR_SERVICE_NAME=app
MONITOR_TRACER_OTLP_ENABLED=false
MONITOR_TRACER_OTLP_ENDPOINT=http://127.0.0.1:4318
```

### YAML 配置文件

详细配置请参考 `configs/config.yaml` 文件。

## 🏗️ 架构设计

本项目采用领域驱动设计（DDD）和清洁架构：

### 分层架构

1. **接口层 (Interfaces)**: 处理外部请求（HTTP、gRPC、WebSocket）
2. **应用层 (Application)**: 业务用例和应用服务
3. **领域层 (Domain)**: 核心业务逻辑和领域模型
4. **基础设施层 (Infrastructure)**: 数据持久化、外部服务集成

### 依赖注入

使用 Google Wire 进行依赖注入，配置文件位于 `internal/di/`：

- `wire.go`: Wire 配置
- `app.go`: 应用程序组装
- `provider.go`: 服务提供者
- `module.go`: 模块定义

## 📊 可观测性

### 链路追踪
- 支持 OpenTelemetry
- 可配置 OTLP 导出器
- 自动追踪 HTTP、gRPC、数据库请求

### 指标监控
- 应用性能指标
- 业务指标
- 系统资源指标

### 日志
- 结构化日志（JSON 格式）
- 日志轮转
- 可配置日志级别

## 🔧 开发工具

### 热重载
```bash
make run  # 使用 fresh 工具进行热重载
```

### 代码生成
```bash
make wire      # 生成依赖注入代码
make wire-show # 查看依赖关系图
```

### 代码质量
```bash
make lint  # 运行 golangci-lint
make test  # 运行测试
```

### Git Hooks
项目包含预配置的 Git hooks：
- `pre-commit`: 提交前代码检查
- `commit-msg`: 提交信息格式检查

## 🚀 部署

### Docker 部署

项目包含 `.dockerignore` 文件，可以构建 Docker 镜像：

```bash
# 构建镜像
docker build -t app-service .

# 运行容器
docker run -p 8080:8080 app-service
```

### 生产环境

1. 设置环境变量 `APP_ENV=prod`
2. 配置生产数据库和 Redis
3. 启用可观测性监控
4. 配置负载均衡和反向代理

## 📝 开发指南

### 添加新功能

1. **创建领域模型** (`internal/domain/`)
2. **实现应用服务** (`internal/application/`)
3. **添加基础设施** (`internal/infrastructure/`)
4. **创建接口层** (`internal/interfaces/`)
5. **配置依赖注入** (`internal/di/`)

### 最佳实践

- 遵循 SOLID 原则
- 使用接口进行解耦
- 编写单元测试
- 保持代码简洁
- 添加适当的日志和监控

## 🤝 贡献

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 支持

如有问题或建议，请：

1. 查看文档
2. 搜索已有 Issues
3. 创建新的 Issue
4. 联系维护者

---

**Happy Coding! 🎉**