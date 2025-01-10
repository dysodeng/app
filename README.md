golang project template
==========
golang项目基础模板

- 包含HTTP、gRPC、WebSocket、消息队列等服务端
- 采用 DDD（领域驱动设计）分层模式，以提高代码的可维护性和可扩展性。

Installation
------------
安装步骤
```sh
git clone https://github.com/dysodeng/app.git
cd app
go mod tidy
cp .env.example .env
```

Usage
-----
使用以下命令启动服务：
```sh
go run main.go
```
