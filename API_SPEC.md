# Todo List API 接口规范

## Base URL
```
http://localhost:8080/api
```

## 数据模型

### Todo Item
```json
{
  "id": 1,
  "title": "完成项目文档",
  "description": "编写 API 设计文档",
  "completed": false,
  "created_at": "2026-02-28T10:30:00Z",
  "updated_at": "2026-02-28T10:30:00Z"
}
```

## API 端点

### 1. 获取所有 Todo
```
GET /todos
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "完成项目文档",
      "description": "编写 API 设计文档",
      "completed": false,
      "created_at": "2026-02-28T10:30:00Z",
      "updated_at": "2026-02-28T10:30:00Z"
    }
  ]
}
```

### 2. 获取单个 Todo
```
GET /todos/:id
```

**路径参数：**
- `id`: Todo ID (整数)

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "完成项目文档",
    "description": "编写 API 设计文档",
    "completed": false,
    "created_at": "2026-02-28T10:30:00Z",
    "updated_at": "2026-02-28T10:30:00Z"
  }
}
```

### 3. 创建 Todo
```
POST /todos
```

**请求体：**
```json
{
  "title": "完成项目文档",
  "description": "编写 API 设计文档"
}
```

**字段验证：**
- `title`: 必填，1-200 字符
- `description`: 可选，最大 1000 字符

**响应示例：**
```json
{
  "code": 0,
  "message": "created",
  "data": {
    "id": 1,
    "title": "完成项目文档",
    "description": "编写 API 设计文档",
    "completed": false,
    "created_at": "2026-02-28T10:30:00Z",
    "updated_at": "2026-02-28T10:30:00Z"
  }
}
```

### 4. 更新 Todo
```
PUT /todos/:id
```

**路径参数：**
- `id`: Todo ID (整数)

**请求体：**
```json
{
  "title": "完成项目文档（更新）",
  "description": "编写 API 设计文档和测试",
  "completed": true
}
```

**字段验证：**
- `title`: 可选，1-200 字符
- `description`: 可选，最大 1000 字符
- `completed`: 可选，布尔值

**响应示例：**
```json
{
  "code": 0,
  "message": "updated",
  "data": {
    "id": 1,
    "title": "完成项目文档（更新）",
    "description": "编写 API 设计文档和测试",
    "completed": true,
    "created_at": "2026-02-28T10:30:00Z",
    "updated_at": "2026-02-28T11:00:00Z"
  }
}
```

### 5. 删除 Todo
```
DELETE /todos/:id
```

**路径参数：**
- `id`: Todo ID (整数)

**响应示例：**
```json
{
  "code": 0,
  "message": "deleted",
  "data": null
}
```

## 错误响应格式

### 通用错误响应
```json
{
  "code": 1,
  "message": "error description",
  "data": null
}
```

### HTTP 状态码
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `404`: 资源不存在
- `500`: 服务器内部错误

### 错误码
- `0`: 成功
- `1001`: 参数验证失败
- `1002`: 资源不存在
- `1003`: 数据库错误
- `1004`: 内部服务错误

## CORS 配置
允许前端跨域访问：
```
Allow-Origin: *
Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Allow-Headers: Content-Type, Authorization
```

## 数据库 Schema (SQLite)

### todos 表
```sql
CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_todos_completed ON todos(completed);
CREATE INDEX idx_todos_created_at ON todos(created_at);
```
