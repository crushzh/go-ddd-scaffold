# Go DDD Scaffold

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrzhoong/go-ddd-scaffold?style=flat-square)](https://goreportcard.com/report/github.com/mrzhoong/go-ddd-scaffold)

> 基于领域驱动设计（DDD）四层架构的 Go 项目脚手架，集成 Gin、GORM、JWT 认证、DI 容器、Swagger 文档和代码生成器。

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

## 特性

- **DDD 四层架构** — interfaces -> application -> domain <- infrastructure
- **依赖反转** — 领域层零外部依赖，基础设施层实现领域接口
- **DI 容器** — 统一的依赖注入管理
- **Gin** HTTP 框架，内置 Recovery、CORS、请求 ID、日志、超时中间件
- **GORM** ORM，支持 SQLite / MySQL / PostgreSQL
- **JWT** 认证，支持角色权限控制
- **Swagger** API 文档自动生成
- **代码生成器** — 一条命令生成完整 DDD CRUD 模块（7 个文件）
- **跨平台编译** — Linux (amd64/arm64/arm32)、Windows、macOS
- **Docker** 多阶段构建 + docker-compose
- **前端嵌入** — `go:embed` 内嵌 SPA 前端
- **结构化日志** — Zap + Lumberjack 日志轮转
- **统一响应格式** — 标准错误码体系
- **优雅退出** — 信号处理
- **服务管理脚本** — systemd / 守护进程

## 架构

```
interfaces -> application -> domain <- infrastructure
```

```
┌────────────┐    ┌─────────────┐    ┌──────────┐    ┌────────────────┐
│ Interfaces │───>│ Application │───>│  Domain  │<───│ Infrastructure │
│  (HTTP)    │    │  (编排)      │    │  (核心)   │    │  (持久化)       │
└────────────┘    └─────────────┘    └──────────┘    └────────────────┘
```

| 层 | 职责 | 依赖 |
|---|------|------|
| **interfaces** | HTTP 请求处理、参数校验、响应转换 | application |
| **application** | 业务编排、事务管理、DTO 转换 | domain |
| **domain** | 核心业务逻辑、实体、值对象、仓储接口 | **无依赖** |
| **infrastructure** | 数据库实现、外部服务、缓存 | domain（实现接口） |

**核心原则**：
- 领域层 **不依赖** 任何外层
- 基础设施层 **实现** 领域层定义的接口（依赖反转）
- 应用层 **编排** 领域逻辑，不包含业务规则

## 技术栈

| 库 | 版本 | 用途 |
|---|------|------|
| [Go](https://go.dev/) | 1.21+ | 语言 |
| [Gin](https://gin-gonic.com/) | v1.9 | HTTP 框架 |
| [GORM](https://gorm.io/) | v1.25 | ORM（SQLite/MySQL/PostgreSQL） |
| [Viper](https://github.com/spf13/viper) | v1.18 | 配置管理 |
| [Zap](https://github.com/uber-go/zap) | v1.26 | 结构化日志 |
| [JWT](https://github.com/golang-jwt/jwt) | v5 | 认证 |
| [Swag](https://github.com/swaggo/swag) | v1.16 | Swagger 文档生成 |
| [Lumberjack](https://github.com/natefinish/lumberjack) | v2.0 | 日志轮转 |

## 环境要求

- Go 1.21+
- （可选）[swag](https://github.com/swaggo/swag)：`go install github.com/swaggo/swag/cmd/swag@latest`

## 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/mrzhoong/go-ddd-scaffold.git
cd go-ddd-scaffold

# 2. 安装依赖
go mod download

# 3. 运行
make run
# 服务启动: http://localhost:8080

# 4. 测试
curl http://localhost:8080/health

# 5. 登录（默认账号: admin / admin123）
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## 项目结构

```
go-ddd-scaffold/
├── cmd/
│   ├── server/main.go                          # 程序入口
│   └── gen/main.go                             # DDD 代码生成器
├── internal/
│   ├── interfaces/                             # 【接口层】
│   │   └── http/
│   │       ├── handler/                        #   HTTP 处理器
│   │       ├── middleware/                      #   中间件
│   │       └── router/                         #   路由注册
│   ├── application/                            # 【应用层】
│   │   ├── service/                            #   应用服务（编排逻辑）
│   │   └── dto/                                #   数据传输对象
│   ├── domain/                                 # 【领域层】★ 核心
│   │   └── example/
│   │       ├── entity.go                       #   聚合根 / 实体
│   │       └── repository.go                   #   仓储接口
│   ├── infrastructure/                         # 【基础设施层】
│   │   └── persistence/database/
│   │       ├── db.go                           #   数据库连接
│   │       ├── example_model.go                #   GORM 模型
│   │       └── example_repo.go                 #   仓储实现
│   └── container/                              # DI 依赖注入容器
│       └── container.go
├── pkg/                                        # 公共工具包
│   ├── config/                                 #   Viper 配置
│   ├── logger/                                 #   Zap 日志
│   └── response/                               #   统一响应
├── templates/                                  # 代码生成模板（7 个）
├── configs/config.yaml                         # 配置文件
├── scripts/                                    # 部署脚本
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## 代码生成器

一条命令生成完整 DDD CRUD 模块：

```bash
make gen name=order cn=订单
```

生成 **7 个文件** 并自动注册路由 + 依赖注入：

| 文件 | 层 | 说明 |
|------|---|------|
| `internal/domain/order/entity.go` | 领域 | 领域实体 + 业务方法 |
| `internal/domain/order/repository.go` | 领域 | 仓储接口 |
| `internal/infrastructure/.../order_model.go` | 基础设施 | GORM 数据模型 |
| `internal/infrastructure/.../order_repo.go` | 基础设施 | 仓储实现 |
| `internal/application/dto/order_dto.go` | 应用 | 数据传输对象 |
| `internal/application/service/order_service.go` | 应用 | 应用服务 |
| `internal/interfaces/http/handler/order_handler.go` | 接口 | HTTP CRUD 处理器 + Swagger |

自动注册：
- `router.go` — 路由注册
- `container.go` — 服务 + 迁移注册

## 配置

通过 `configs/config.yaml` 加载配置，支持 `APP_` 前缀的环境变量覆盖。

```yaml
app:
  name: "myapp"
  mode: "debug"           # debug, release, test

server:
  host: "0.0.0.0"
  port: 8080

database:
  type: "sqlite"          # sqlite, mysql, postgres
  path: "./data/app.db"
  # host: "127.0.0.1"    # MySQL/PostgreSQL
  # port: 3306
  # username: "root"
  # password: ""
  # dbname: "mydb"

jwt:
  secret: "change-me-in-production"
  expire: 24              # 小时
  refresh_hours: 168      # 7 天

log:
  level: "info"           # debug, info, warn, error
  filename: "logs/app.log"
  max_size: 100           # MB
  max_backups: 10
  max_age: 30             # 天
```

## API 示例

```bash
# 健康检查
curl http://localhost:8080/health

# 登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# 创建
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"hello"}'

# 列表（分页）
curl "http://localhost:8080/api/v1/examples?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 查询
curl http://localhost:8080/api/v1/examples/1 \
  -H "Authorization: Bearer $TOKEN"

# 更新
curl -X PUT http://localhost:8080/api/v1/examples/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"updated","description":"world"}'

# 删除
curl -X DELETE http://localhost:8080/api/v1/examples/1 \
  -H "Authorization: Bearer $TOKEN"
```

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| 错误码范围 | 分类 | 说明 |
|-----------|------|------|
| 0 | 成功 | 操作成功 |
| 1001-1999 | 客户端 | 参数/校验错误 |
| 2001-2999 | 资源 | 不存在、冲突 |
| 3001-3999 | 业务 | 业务逻辑错误 |
| 4001-4999 | 认证 | 未认证、无权限、Token 过期 |
| 5001-5999 | 系统 | 内部错误、数据库错误、超时 |

## 添加新模块

### 自动生成（推荐）

```bash
make gen name=order cn=订单
# 然后：
# 1. 编辑 internal/domain/order/entity.go         — 添加领域字段和业务方法
# 2. 编辑 internal/infrastructure/.../order_model.go — 同步数据库字段
# 3. 编辑 internal/application/dto/order_dto.go      — 同步 API 字段
# 4. 执行 make docs                                   — 更新 Swagger 文档
```

### 手动添加

1. **领域层**：`internal/domain/<module>/` — 定义实体、值对象、仓储接口
2. **基础设施层**：`internal/infrastructure/persistence/database/` — 实现仓储
3. **应用层**：`internal/application/service/` — 创建应用服务；`dto/` — 定义 DTO
4. **接口层**：`internal/interfaces/http/handler/` — 创建 HTTP 处理器
5. **容器**：`internal/container/container.go` — 注册依赖
6. **路由**：`internal/interfaces/http/router/router.go` — 注册路由

## 部署

### Docker

```bash
docker-compose up -d
```

### 跨平台编译

```bash
make build-all            # 全平台编译
make build-linux          # Linux amd64
make build-arm64          # Linux arm64
make build-arm32          # Linux arm32
make build-windows        # Windows
```

### 服务管理

```bash
./scripts/manage.sh start     # 启动（含守护进程）
./scripts/manage.sh stop      # 优雅停止
./scripts/manage.sh restart   # 重启
./scripts/manage.sh status    # 查看状态
```

## 常用命令

| 命令 | 说明 |
|------|------|
| `make run` | 开发运行 |
| `make build` | 本地编译 |
| `make build-all` | 全平台编译 |
| `make gen name=order cn=订单` | 生成 DDD 模块 |
| `make docs` | 生成 Swagger 文档 |
| `make test` | 运行测试 |
| `make lint` | 运行代码检查 |
| `make clean` | 清理构建产物 |
| `make help` | 查看所有命令 |

## 自定义

重命名模块：

```bash
# macOS
find . -name "*.go" -exec sed -i '' 's|go-ddd-scaffold|my-project|g' {} +
sed -i '' 's|go-ddd-scaffold|my-project|g' go.mod Makefile configs/config.yaml

# Linux
find . -name "*.go" -exec sed -i 's|go-ddd-scaffold|my-project|g' {} +
sed -i 's|go-ddd-scaffold|my-project|g' go.mod Makefile configs/config.yaml

go mod tidy
```

## 贡献

参见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可证

[MIT](LICENSE)
