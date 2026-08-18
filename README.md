# 设备巡检与维护平台

这是一个从 0 到 1 的完整前后端项目，面向设备管理人员和现场巡检人员，支持设备台账、巡检计划、巡检任务、巡检记录、异常处理和维修记录的闭环管理。后端使用 Go，前端使用 React + TypeScript + Vite。

## 功能概览

- 设备台账：设备编号唯一，支持按关键词、类别、状态、负责人、安装位置筛选，可查看设备详情、巡检历史和维修历史。
- 巡检计划：支持指定设备或设备类别，周期可选择每日、每周、每月，配置开始/结束时间、负责人和多个检查项目；计划可启用、暂停、终止。
- 巡检任务：后台按计划自动生成任务，同一设备同一周期内去重；任务支持草稿、正常完成、发现异常，必填项校验，已完成任务不可重复提交。
- 异常管理：巡检发现异常时自动创建异常记录，支持指派、进展更新、原因分析、关闭说明和验证结果。
- 维修记录：记录维修类型、内容、单位、起止时间、费用和结果，已关闭异常可关联维修记录。
- 用户与登录：账号密码登录，bcrypt 加密存储密码，HMAC 短期访问令牌；当前不设置角色分级。
- 服务治理：结构化日志、请求日志、请求 ID、panic 恢复、超时、优雅退出、`/healthz`、`/readyz`、`/metrics`。

## 技术栈

- 后端：Go 1.24+、标准库 `net/http`、`modernc.org/sqlite`（开发默认）、`github.com/lib/pq`（生产 PostgreSQL）、`golang.org/x/crypto/bcrypt`、`gopkg.in/yaml.v3`。
- 前端：React 18、TypeScript、Vite、React Router、lucide-react。

## 目录结构

```text
.
├── api/                     # REST 路由与 OpenAPI 文档
│   ├── router.go
│   └── openapi.yaml
├── cmd/server/main.go       # 服务入口、依赖注入、优雅退出
├── configs/
│   ├── config.yaml          # 本地开发配置（SQLite）
│   └── config.prod.yaml     # 生产配置（PostgreSQL）
├── internal/
│   ├── domain/              # 共享领域模型、仓储接口、错误
│   ├── config/              # YAML + 环境变量配置
│   ├── platform/            # 数据库、认证、中间件、Web 工具
│   ├── device/              # 设备台账领域
│   ├── plan/                # 巡检计划领域
│   ├── task/                # 巡检任务领域
│   ├── record/              # 巡检记录领域
│   ├── anomaly/             # 异常管理领域
│   ├── repair/              # 维修记录领域
│   ├── user/                # 用户与登录领域
│   └── scheduler/           # 周期任务生成器
├── migrations/              # PostgreSQL 迁移脚本
├── deploy/                  # Dockerfile 与 Compose
├── scripts/
│   ├── run-dev.sh           # 本地同时启动前后端
│   └── migrate.sh           # 执行 PostgreSQL 迁移
└── frontend/                # React + TypeScript + Vite
    └── src/
        ├── api/             # 请求封装与服务 API
        ├── components/      # 通用组件
        ├── pages/           # 页面
        ├── types/           # 类型定义
        └── utils/           # 工具函数
```

每个业务领域内部按职责拆分为：

- `application`：业务服务与用例逻辑。
- `adapter/http`：HTTP 参数解析与响应。
- `infrastructure`：SQLite 数据访问实现。

## 快速启动

环境要求：Go 约 1.24+、Node 约 22、npm 可用。

```bash
cd inspection-platform
chmod +x scripts/run-dev.sh
./scripts/run-dev.sh
```

该脚本会：

1. 使用 `configs/config.yaml` 启动 Go 后端，默认地址 `http://localhost:8080`。
2. 启动 Vite 前端开发服务器，默认地址 `http://localhost:5173`。
3. 后端首次启动时自动创建 SQLite 表并初始化默认账号。

默认账号：

- 用户名：`admin`
- 密码：`admin123`

前端页面可通过 Vite 代理访问 `/api`、`/healthz`、`/readyz` 和 `/metrics`。

## 本地验证

仅使用 Go、Node 和 curl 完成验证，不依赖外部浏览器或 Chrome。

健康检查：

```bash
curl -sS -i http://localhost:8080/healthz
curl -sS -i http://localhost:8080/readyz
```

业务闭环验证示例：

```bash
TOKEN=$(curl -sS -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | node -e "let s='';process.stdin.on('data',d=>s+=d);process.stdin.on('end',()=>console.log(JSON.parse(s).token))")

curl -sS -X POST http://localhost:8080/api/devices \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"code":"DEV-001","name":"空压机A","category":"机械设备","model":"AC-200","serialNumber":"SN-AC200-01","manufacturer":"示例厂商","location":"动力站房","installDate":"2026-01-10","warrantyEnd":"2028-01-10","owner":"王工","status":"active","remark":"主用设备"}'
```

完整验证链路包括：登录、创建设备、创建计划并生成任务、查询任务、提交异常巡检结果、关闭异常、创建关联维修记录、查询设备历史。

## API 说明

除 `POST /api/auth/login`、`/healthz`、`/readyz`、`/metrics` 外，其余接口都需要 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/auth/login` | 登录并返回短期访问令牌 |
| GET | `/api/auth/me` | 当前登录用户 |
| GET | `/api/users` | 用户列表 |
| GET | `/api/dictionaries` | 类别、位置和状态字典 |
| GET | `/api/devices` | 分页筛选设备 |
| POST | `/api/devices` | 创建设备 |
| GET | `/api/devices/{id}` | 设备详情 |
| PUT | `/api/devices/{id}` | 更新设备 |
| GET | `/api/devices/{id}/history` | 设备巡检与维修历史 |
| GET | `/api/plans` | 分页筛选计划 |
| POST | `/api/plans` | 创建计划并生成近期任务 |
| GET | `/api/plans/{id}` | 计划详情 |
| PUT | `/api/plans/{id}` | 更新计划 |
| POST | `/api/plans/{id}/status` | 启用、暂停或终止计划 |
| POST | `/api/plans/generate` | 立即生成一次未来任务 |
| GET | `/api/tasks` | 分页筛选任务 |
| GET | `/api/tasks/{id}` | 任务详情 |
| POST | `/api/tasks/{id}/draft` | 保存草稿 |
| POST | `/api/tasks/{id}/submit` | 提交巡检结果 |
| GET | `/api/records?deviceId=` | 查询设备巡检记录 |
| GET | `/api/anomalies` | 分页筛选异常 |
| GET | `/api/anomalies/{id}` | 异常详情 |
| POST | `/api/anomalies/{id}/assign` | 指派负责人 |
| POST | `/api/anomalies/{id}/progress` | 更新处理进展 |
| POST | `/api/anomalies/{id}/close` | 关闭异常并验证 |
| GET | `/api/repairs` | 分页筛选维修记录 |
| POST | `/api/repairs` | 创建维修记录 |
| GET | `/api/repairs/{id}` | 维修记录详情 |
| PUT | `/api/repairs/{id}` | 更新维修记录 |

更完整的接口描述见 `api/openapi.yaml`。

## 配置项

配置通过 YAML 文件加载，并可用环境变量覆盖部分关键项：

| 配置项 | 说明 | 环境变量 |
| --- | --- | --- |
| `server.port` | HTTP 监听端口 | `APP_PORT` |
| `server.read_timeout` | 读超时 | 无 |
| `server.write_timeout` | 写超时 | 无 |
| `database.driver` | `sqlite` 或 `postgres` | `DB_DRIVER` |
| `database.sqlite_path` | SQLite 文件路径 | `SQLITE_PATH` |
| `database.postgres_dsn` | PostgreSQL DSN | `POSTGRES_DSN` |
| `auth.jwt_secret` | 访问令牌签名密钥 | `JWT_SECRET` |
| `auth.token_ttl` | 访问令牌有效期 | 无 |
| `scheduler.enabled` | 是否启用后台生成任务 | 无 |
| `scheduler.interval` | 生成周期 | 无 |
| `log.level` | 日志级别 | `LOG_LEVEL` |
| `log.format` | `json` 或 `text` | 无 |

## 数据库与迁移

### 本地 SQLite

默认 `database.driver: sqlite`，数据库文件位于 `data/inspection.db`。服务启动时会自动执行 `internal/platform/database/schema.go` 中的 DDL，无需外部数据库服务。

### 生产 PostgreSQL

1. 启动 PostgreSQL。
2. 使用 `configs/config.prod.yaml`，将 `POSTGRES_DSN` 指向目标数据库。
3. 执行迁移：

```bash
export POSTGRES_DSN='postgres://inspection:inspection@localhost:5432/inspection?sslmode=disable'
./scripts/migrate.sh
```

迁移文件：

- `migrations/0001_init.sql`：创建全部业务表。
- `migrations/0002_seed.sql`：创建默认 `admin` 账号。

SQLite 用于本地开发验证，PostgreSQL 迁移用于生产环境；两者字段和约束保持一致。

## 构建与测试结果

在项目根目录执行：

```bash
gofmt -w $(find cmd api internal -name '*.go' -type f)
go mod tidy
go vet ./...
go build ./...
go test ./...
cd frontend && npm install && npm run build
```

结果：

- `go vet ./...`：通过，无报告。
- `go build ./...`：通过，无错误。
- `go test ./...`：通过，当前仓库没有测试文件。
- `npm run build`：通过，Vite 输出 `dist/`。

### 源码行数

扣除 `*_test.go`、`node_modules`、lockfile 和构建产物后：

```bash
GO_LINES=$(find cmd api internal -name '*.go' -type f ! -name '*_test.go' -print0 | xargs -0 wc -l | tail -1 | awk '{print $1}')
TS_LINES=$(find frontend -type f \( -name '*.ts' -o -name '*.tsx' \) -not -path '*/node_modules/*' -not -path '*/dist/*' -print0 | xargs -0 wc -l | tail -1 | awk '{print $1}')
echo "Go: $GO_LINES"
echo "TypeScript/TSX: $TS_LINES"
echo "Total: $((GO_LINES + TS_LINES))"
```

当前统计：

- Go 源码：4150 行
- TypeScript/TSX 源码：1662 行
- 总计：5812 行
