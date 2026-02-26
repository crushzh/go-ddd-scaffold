# Go DDD Scaffold

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/crushzh/go-ddd-scaffold?style=flat-square)](https://goreportcard.com/report/github.com/crushzh/go-ddd-scaffold)

> 基於領域驅動設計（DDD）四層架構的 Go 專案腳手架，整合 Gin、GORM、JWT 認證、DI 容器、Swagger 文件和程式碼產生器。

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md)

## 特色

- **DDD 四層架構** — interfaces -> application -> domain <- infrastructure
- **依賴反轉** — 領域層零外部依賴，基礎設施層實現領域介面
- **DI 容器** — 統一的依賴注入管理
- **Gin** HTTP 框架，內建 Recovery、CORS、請求 ID、日誌、逾時中介軟體
- **GORM** ORM，支援 SQLite / MySQL / PostgreSQL
- **JWT** 認證，支援角色權限控制
- **Swagger** API 文件自動產生
- **程式碼產生器** — 一條指令產生完整 DDD CRUD 模組（7 個檔案）
- **跨平台編譯** — Linux (amd64/arm64/arm32)、Windows、macOS
- **Docker** 多階段建置 + docker-compose
- **前端嵌入** — `go:embed` 內嵌 SPA 前端
- **前端範本** — UmiJS Max + Ant Design ProComponents（登入 + 儀表板 + CRUD 範例）
- **Swagger UI** — 預整合，啟動即可存取 `/swagger/index.html`
- **SQLite 自動初始化** — 自動建立資料目錄 + 預設管理員種子資料（admin/admin123）
- **自解壓安裝包** — `.run` 一鍵安裝包建置（`make package-all`）
- **結構化日誌** — Zap + Lumberjack 日誌輪替
- **統一回應格式** — 標準錯誤碼體系
- **優雅退出** — 訊號處理
- **服務管理腳本** — systemd / 守護程序 + 解除安裝腳本

## 架構

```
interfaces -> application -> domain <- infrastructure
```

```
┌────────────┐    ┌─────────────┐    ┌──────────┐    ┌────────────────┐
│ Interfaces │───>│ Application │───>│  Domain  │<───│ Infrastructure │
│  (HTTP)    │    │  (編排)      │    │  (核心)   │    │  (持久化)       │
└────────────┘    └─────────────┘    └──────────┘    └────────────────┘
```

| 層 | 職責 | 依賴 |
|---|------|------|
| **interfaces** | HTTP 請求處理、參數驗證、回應轉換 | application |
| **application** | 業務編排、交易管理、DTO 轉換 | domain |
| **domain** | 核心業務邏輯、實體、值物件、儲存庫介面 | **無依賴** |
| **infrastructure** | 資料庫實現、外部服務、快取 | domain（實現介面） |

**核心原則**：
- 領域層 **不依賴** 任何外層
- 基礎設施層 **實現** 領域層定義的介面（依賴反轉）
- 應用層 **編排** 領域邏輯，不包含業務規則

## 技術棧

| 函式庫 | 版本 | 用途 |
|-------|------|------|
| [Go](https://go.dev/) | 1.21+ | 語言 |
| [Gin](https://gin-gonic.com/) | v1.9 | HTTP 框架 |
| [GORM](https://gorm.io/) | v1.25 | ORM（SQLite/MySQL/PostgreSQL） |
| [Viper](https://github.com/spf13/viper) | v1.18 | 設定管理 |
| [Zap](https://github.com/uber-go/zap) | v1.26 | 結構化日誌 |
| [JWT](https://github.com/golang-jwt/jwt) | v5 | 認證 |
| [Swag](https://github.com/swaggo/swag) | v1.16 | Swagger 文件產生 |
| [Lumberjack](https://github.com/natefinish/lumberjack) | v2.0 | 日誌輪替 |

## 環境需求

- Go 1.21+
- （選用）[swag](https://github.com/swaggo/swag)：`go install github.com/swaggo/swag/cmd/swag@latest`

## 快速開始

```bash
# 1. 複製儲存庫
git clone https://github.com/crushzh/go-ddd-scaffold.git
cd go-ddd-scaffold

# 2. 安裝依賴
go mod download

# 3. 執行
make run
# 服務啟動: http://localhost:8080

# 4. 測試
curl http://localhost:8080/health

# 5. 登入（預設帳號: admin / admin123）
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## 專案結構

```
go-ddd-scaffold/
├── cmd/
│   ├── server/main.go                          # 程式進入點
│   └── gen/main.go                             # DDD 程式碼產生器
├── internal/
│   ├── interfaces/                             # 【介面層】
│   │   └── http/
│   │       ├── handler/                        #   HTTP 處理器
│   │       ├── middleware/                      #   中介軟體
│   │       └── router/                         #   路由註冊
│   ├── application/                            # 【應用層】
│   │   ├── service/                            #   應用服務（編排邏輯）
│   │   └── dto/                                #   資料傳輸物件
│   ├── domain/                                 # 【領域層】★ 核心
│   │   └── example/
│   │       ├── entity.go                       #   聚合根 / 實體
│   │       └── repository.go                   #   儲存庫介面
│   ├── infrastructure/                         # 【基礎設施層】
│   │   └── persistence/database/
│   │       ├── db.go                           #   資料庫連線
│   │       ├── example_model.go                #   GORM 模型
│   │       └── example_repo.go                 #   儲存庫實現
│   └── container/                              # DI 依賴注入容器
│       └── container.go
├── pkg/                                        # 公共工具套件
│   ├── config/                                 #   Viper 設定
│   ├── logger/                                 #   Zap 日誌
│   └── response/                               #   統一回應
├── web/                                        # 前端（UmiJS Max + ProComponents）
├── docs/swagger/                               # Swagger 文件（預產生）
├── templates/                                  # 程式碼產生範本（7 個）
├── configs/config.yaml                         # 設定檔
├── scripts/                                    # 部署腳本（安裝/解除安裝/管理）
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## 程式碼產生器

一條指令產生完整 DDD CRUD 模組：

```bash
make gen name=order cn=訂單
```

產生 **7 個檔案** 並自動註冊路由 + 依賴注入：

| 檔案 | 層 | 說明 |
|------|---|------|
| `internal/domain/order/entity.go` | 領域 | 領域實體 + 業務方法 |
| `internal/domain/order/repository.go` | 領域 | 儲存庫介面 |
| `internal/infrastructure/.../order_model.go` | 基礎設施 | GORM 資料模型 |
| `internal/infrastructure/.../order_repo.go` | 基礎設施 | 儲存庫實現 |
| `internal/application/dto/order_dto.go` | 應用 | 資料傳輸物件 |
| `internal/application/service/order_service.go` | 應用 | 應用服務 |
| `internal/interfaces/http/handler/order_handler.go` | 介面 | HTTP CRUD 處理器 + Swagger |

自動註冊：
- `router.go` — 路由註冊
- `container.go` — 服務 + 遷移註冊

## 設定

透過 `configs/config.yaml` 載入設定，支援 `APP_` 前綴的環境變數覆寫。

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

jwt:
  secret: "change-me-in-production"
  expire: 24              # 小時
  refresh_hours: 168      # 7 天

log:
  level: "info"           # debug, info, warn, error
  filename: "logs/app.log"
  max_size: 100           # MB
  max_backups: 10
  max_age: 30             # 天
```

## API 範例

```bash
# 健康檢查
curl http://localhost:8080/health

# 登入
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.data.token')

# 建立
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test","description":"hello"}'

# 列表（分頁）
curl "http://localhost:8080/api/v1/examples?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

## 統一回應格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

| 錯誤碼範圍 | 分類 | 說明 |
|-----------|------|------|
| 0 | 成功 | 操作成功 |
| 1001-1999 | 客戶端 | 參數/驗證錯誤 |
| 2001-2999 | 資源 | 不存在、衝突 |
| 3001-3999 | 業務 | 業務邏輯錯誤 |
| 4001-4999 | 認證 | 未認證、無權限、Token 過期 |
| 5001-5999 | 系統 | 內部錯誤、資料庫錯誤、逾時 |

## 新增模組

### 自動產生（建議）

```bash
make gen name=order cn=訂單
# 然後：
# 1. 編輯 internal/domain/order/entity.go         — 新增領域欄位和業務方法
# 2. 編輯 internal/infrastructure/.../order_model.go — 同步資料庫欄位
# 3. 編輯 internal/application/dto/order_dto.go      — 同步 API 欄位
# 4. 執行 make docs                                   — 更新 Swagger 文件
```

### 手動新增

1. **領域層**：`internal/domain/<module>/` — 定義實體、值物件、儲存庫介面
2. **基礎設施層**：`internal/infrastructure/persistence/database/` — 實現儲存庫
3. **應用層**：`internal/application/service/` — 建立應用服務；`dto/` — 定義 DTO
4. **介面層**：`internal/interfaces/http/handler/` — 建立 HTTP 處理器
5. **容器**：`internal/container/container.go` — 註冊依賴
6. **路由**：`internal/interfaces/http/router/router.go` — 註冊路由

## 部署

### Docker

```bash
docker-compose up -d
```

### 跨平台編譯

```bash
make build-all            # 全平台編譯
make build-linux          # Linux amd64
make build-arm64          # Linux arm64
make build-arm32          # Linux arm32
make build-windows        # Windows
```

### 服務管理

```bash
./scripts/manage.sh start     # 啟動（含守護程序）
./scripts/manage.sh stop      # 優雅停止
./scripts/manage.sh restart   # 重新啟動
./scripts/manage.sh status    # 檢視狀態
```

## 前端

```bash
cd web
npm install
npm run dev     # 開發模式（http://localhost:8000，代理到 :8080）
npm run build   # 生產建置（輸出到 ../internal/web/dist/）
```

建置後的前端透過 `go:embed` 嵌入 Go 二進位檔。執行 `make web && make build` 即可產生包含前端的單一可執行檔。

## 自訂

### 使用 init.sh（建議）

如果從 [go-scaffold](https://github.com/crushzh/go-scaffold) 複製，使用初始化腳本：

```bash
./init.sh ddd my-service
# 自動完成: 複製範本 → 替換模組名 → go mod tidy → swag init → git init
```

### 手動重新命名

```bash
# macOS
find . -name "*.go" -exec sed -i '' 's|go-ddd-scaffold|my-project|g' {} +
sed -i '' 's|go-ddd-scaffold|my-project|g' go.mod Makefile configs/config.yaml
go mod tidy
```

## 貢獻

參見 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 授權條款

[MIT](LICENSE)
