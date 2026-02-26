# Go DDD Scaffold

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/crushzh/go-ddd-scaffold?style=flat-square)](https://goreportcard.com/report/github.com/crushzh/go-ddd-scaffold)

> Production-ready Go project scaffold based on Domain-Driven Design (DDD) with clean 4-layer architecture, Gin, GORM, JWT auth, DI container, Swagger docs, and code generator.

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

## Features

- **DDD 4-Layer Architecture** — interfaces -> application -> domain <- infrastructure
- **Dependency Inversion** — Domain layer has zero external dependencies; infrastructure implements domain interfaces
- **DI Container** — Centralized dependency injection management
- **Gin** HTTP framework with Recovery, CORS, Request ID, Logging, Timeout middleware
- **GORM** ORM supporting SQLite / MySQL / PostgreSQL
- **JWT** authentication with role-based access control
- **Swagger** API documentation auto-generation
- **Code Generator** — Single command generates full DDD CRUD module (7 files)
- **Cross-platform Build** — Linux (amd64/arm64/arm32), Windows, macOS
- **Docker** multi-stage build + docker-compose
- **Frontend Embedding** — `go:embed` for SPA static files
- **Frontend Template** — UmiJS Max + Ant Design ProComponents (Login + Dashboard + CRUD)
- **Swagger UI** — Pre-integrated, visit `/swagger/index.html` on startup
- **SQLite Auto-setup** — Auto-creates data directory + default admin seed (admin/admin123)
- **Self-extracting Installer** — `.run` package build (`make package-all`)
- **Structured Logging** — Zap + Lumberjack log rotation
- **Unified Response Format** — Standardized error code system
- **Graceful Shutdown** — Signal handling
- **Service Management** — systemd / daemon scripts + uninstall script

## Architecture

```
interfaces -> application -> domain <- infrastructure
```

```
┌────────────┐    ┌─────────────┐    ┌──────────┐    ┌────────────────┐
│ Interfaces │───>│ Application │───>│  Domain  │<───│ Infrastructure │
│  (HTTP)    │    │ (Orchestr.) │    │  (Core)  │    │  (Persistence) │
└────────────┘    └─────────────┘    └──────────┘    └────────────────┘
```

| Layer | Responsibility | Dependencies |
|-------|---------------|-------------|
| **interfaces** | HTTP request handling, parameter validation, response formatting | application |
| **application** | Business orchestration, transaction management, DTO conversion | domain |
| **domain** | Core business logic, entities, value objects, repository interfaces | **none** |
| **infrastructure** | Database implementation, external services, caching | domain (implements interfaces) |

**Core Principles**:
- Domain layer has **no dependencies** on any outer layer
- Infrastructure layer **implements** interfaces defined by the domain layer (Dependency Inversion)
- Application layer **orchestrates** domain logic without containing business rules

## Tech Stack

| Library | Version | Purpose |
|---------|---------|---------|
| [Go](https://go.dev/) | 1.21+ | Language |
| [Gin](https://gin-gonic.com/) | v1.9 | HTTP framework |
| [GORM](https://gorm.io/) | v1.25 | ORM (SQLite/MySQL/PostgreSQL) |
| [Viper](https://github.com/spf13/viper) | v1.18 | Configuration management |
| [Zap](https://github.com/uber-go/zap) | v1.26 | Structured logging |
| [JWT](https://github.com/golang-jwt/jwt) | v5 | Authentication |
| [Swag](https://github.com/swaggo/swag) | v1.16 | Swagger doc generation |
| [Lumberjack](https://github.com/natefinish/lumberjack) | v2.0 | Log rotation |

## Prerequisites

- Go 1.21+
- (Optional) [swag](https://github.com/swaggo/swag) for Swagger docs: `go install github.com/swaggo/swag/cmd/swag@latest`

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/crushzh/go-ddd-scaffold.git
cd go-ddd-scaffold

# 2. Install dependencies
go mod download

# 3. Run
make run
# Server starts at: http://localhost:8080

# 4. Health check
curl http://localhost:8080/health

# 5. Login (default: admin / admin123)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## Project Structure

```
go-ddd-scaffold/
├── cmd/
│   ├── server/main.go                          # Application entry point
│   └── gen/main.go                             # DDD code generator
├── internal/
│   ├── interfaces/                             # [Interface Layer]
│   │   └── http/
│   │       ├── handler/                        #   HTTP handlers
│   │       ├── middleware/                      #   Middleware
│   │       └── router/                         #   Route registration
│   ├── application/                            # [Application Layer]
│   │   ├── service/                            #   Application services (orchestration)
│   │   └── dto/                                #   Data Transfer Objects
│   ├── domain/                                 # [Domain Layer] ★ Core
│   │   └── example/
│   │       ├── entity.go                       #   Aggregate root / Entity
│   │       └── repository.go                   #   Repository interface
│   ├── infrastructure/                         # [Infrastructure Layer]
│   │   └── persistence/database/
│   │       ├── db.go                           #   Database connection
│   │       ├── example_model.go                #   GORM model
│   │       └── example_repo.go                 #   Repository implementation
│   └── container/                              # DI container
│       └── container.go
├── pkg/                                        # Shared packages
│   ├── config/                                 #   Viper configuration
│   ├── logger/                                 #   Zap logging
│   └── response/                               #   Unified API response
├── web/                                        # Frontend (UmiJS Max + ProComponents)
├── docs/swagger/                               # Swagger docs (pre-generated)
├── templates/                                  # Code generator templates (7 files)
├── configs/config.yaml                         # Configuration file
├── scripts/                                    # Deployment scripts (install/uninstall/manage)
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Code Generator

Generate a full DDD CRUD module with a single command:

```bash
make gen name=order cn=Order
```

This generates **7 files** and auto-registers routes + DI:

| File | Layer | Description |
|------|-------|-------------|
| `internal/domain/order/entity.go` | Domain | Domain entity + business methods |
| `internal/domain/order/repository.go` | Domain | Repository interface |
| `internal/infrastructure/.../order_model.go` | Infrastructure | GORM data model |
| `internal/infrastructure/.../order_repo.go` | Infrastructure | Repository implementation |
| `internal/application/dto/order_dto.go` | Application | Data Transfer Objects |
| `internal/application/service/order_service.go` | Application | Application service |
| `internal/interfaces/http/handler/order_handler.go` | Interface | HTTP CRUD handler + Swagger |

Auto-registration:
- `router.go` — Route registration
- `container.go` — Service + migration registration

## Configuration

Configuration is loaded from `configs/config.yaml` with support for `APP_` prefixed environment variable overrides.

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
  expire: 24              # hours
  refresh_hours: 168      # 7 days

log:
  level: "info"           # debug, info, warn, error
  filename: "logs/app.log"
  max_size: 100           # MB
  max_backups: 10
  max_age: 30             # days
```

## API Examples

```bash
# Health check
curl http://localhost:8080/health

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# Create
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"hello"}'

# List (with pagination)
curl "http://localhost:8080/api/v1/examples?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# Get by ID
curl http://localhost:8080/api/v1/examples/1 \
  -H "Authorization: Bearer $TOKEN"

# Update
curl -X PUT http://localhost:8080/api/v1/examples/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"updated","description":"world"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/examples/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Unified Response Format

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| Code Range | Category | Description |
|-----------|----------|-------------|
| 0 | Success | Operation successful |
| 1001-1999 | Client | Parameter / validation errors |
| 2001-2999 | Resource | Not found, conflict |
| 3001-3999 | Business | Business logic errors |
| 4001-4999 | Auth | Unauthorized, forbidden, token expired |
| 5001-5999 | System | Internal error, database error, timeout |

## Adding a New Module

### Automatic (Recommended)

```bash
make gen name=order cn=Order
# Then:
# 1. Edit internal/domain/order/entity.go   — add domain fields and business methods
# 2. Edit internal/infrastructure/.../order_model.go — sync database fields
# 3. Edit internal/application/dto/order_dto.go      — sync API fields
# 4. Run make docs                                    — update Swagger docs
```

### Manual

1. **Domain**: `internal/domain/<module>/` — define entities, value objects, repository interfaces
2. **Infrastructure**: `internal/infrastructure/persistence/database/` — implement repository
3. **Application**: `internal/application/service/` — create application service; `dto/` — define DTOs
4. **Interfaces**: `internal/interfaces/http/handler/` — create HTTP handler
5. **Container**: `internal/container/container.go` — register dependencies
6. **Router**: `internal/interfaces/http/router/router.go` — register routes

## Deployment

### Docker

```bash
docker-compose up -d
```

### Cross-platform Build

```bash
make build-all            # All platforms
make build-linux          # Linux amd64
make build-arm64          # Linux arm64
make build-arm32          # Linux arm32
make build-windows        # Windows
```

### Service Management

```bash
./scripts/manage.sh start     # Start (with daemon)
./scripts/manage.sh stop      # Graceful stop
./scripts/manage.sh restart   # Restart
./scripts/manage.sh status    # Check status
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make run` | Development run |
| `make build` | Local build |
| `make build-all` | Cross-platform build |
| `make gen name=order cn=Order` | Generate DDD module |
| `make docs` | Generate Swagger docs |
| `make web` | Build frontend (UmiJS) |
| `make package-all` | Build .run installers (all platforms) |
| `make package-linux` | Build .run installer (amd64) |
| `make package-arm64` | Build .run installer (arm64) |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make clean` | Clean build artifacts |
| `make help` | Show all commands |

## Frontend

```bash
cd web
npm install
npm run dev     # Development (http://localhost:8000, proxy to :8080)
npm run build   # Production build (output: ../internal/web/dist/)
```

Built frontend is embedded into the Go binary via `go:embed`. Run `make web && make build` to produce a single binary with frontend included.

## Customization

### Using init.sh (Recommended)

If you cloned this template from [go-scaffold](https://github.com/crushzh/go-scaffold), use the initialization script:

```bash
./init.sh ddd my-service
# Automatically: copy template → replace module name → go mod tidy → swag init → git init
```

### Manual Rename

```bash
# macOS
find . -name "*.go" -exec sed -i '' 's|go-ddd-scaffold|my-project|g' {} +
sed -i '' 's|go-ddd-scaffold|my-project|g' go.mod Makefile configs/config.yaml

# Linux
find . -name "*.go" -exec sed -i 's|go-ddd-scaffold|my-project|g' {} +
sed -i 's|go-ddd-scaffold|my-project|g' go.mod Makefile configs/config.yaml

go mod tidy
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
