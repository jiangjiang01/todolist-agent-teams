# Todo List 后端服务

基于 Go + Gin + SQLite 的 RESTful API 服务。

## 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # 入口文件
├── internal/
│   ├── database/
│   │   └── database.go      # 数据库连接和初始化
│   ├── models/
│   │   └── todo.go          # Todo 数据模型
│   ├── handlers/
│   │   └── todo.go          # API 处理器
│   └── middleware/
│       └── cors.go          # CORS 中间件
├── bin/
│   └── server               # 编译后的二进制文件
├── go.mod
├── go.sum
└── API_SPEC.md              # API 规范文档
```

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 编译

```bash
go build -o bin/server cmd/server/main.go
```

### 运行

```bash
# 使用默认配置（端口 8080，数据库路径 ./data/todos.db）
./bin/server

# 自定义端口
PORT=3000 ./bin/server

# 自定义数据库路径
DB_PATH=/path/to/todos.db ./bin/server
```

### 或直接使用 go run

```bash
go run cmd/server/main.go
```

## API 端点

- `GET /api/todos` - 获取所有 Todo
- `GET /api/todos/:id` - 获取单个 Todo
- `POST /api/todos` - 创建 Todo
- `PUT /api/todos/:id` - 更新 Todo
- `DELETE /api/todos/:id` - 删除 Todo

详细 API 规范请参考 [API_SPEC.md](./API_SPEC.md)

## 技术栈

- **Go 1.21+**
- **Gin** - Web 框架
- **SQLite3** - 数据库
- **CORS** - 跨域支持

## 设计原则

- **单一职责**：每个模块负责单一功能
- **依赖注入**：数据库连接统一管理
- **统一响应**：标准化的 API 响应格式
- **错误处理**：清晰的错误码和消息
