# GFRD Web 管理后台

## 快速开始

### 启动 Web 服务

```bash
cd gfrd/gen

# 方式 1: 使用 web 命令
go run main.go web

# 方式 2: 使用 serve 命令（同上）
go run main.go serve

# 方式 3: 运行编译后的二进制文件
./gfrd web
```

### 访问地址

启动后访问：http://localhost:8199

## 功能特性

### 1. 数据库配置
- 支持 MySQL 和 PostgreSQL
- 测试连接功能
- 连接字符串格式：
  - MySQL: `mysql:root:123456@tcp(127.0.0.1:3306)/gfrd`
  - PostgreSQL: `postgres:root:123456@tcp(127.0.0.1:5432)/gfrd`

### 2. 表选择
- 显示所有表及其注释
- 支持多选
- 支持全选操作

### 3. 字段配置
- 配置列表显示字段
- 配置查询条件字段
- 设置表单类型（input/textarea/select/radio/checkbox/switch/date/datetime）

### 4. 代码生成
- 选择要生成的功能（列表/新增/修改/删除/详情）
- 配置输出目录
- 实时生成代码

### 5. 生成历史
- 查看所有生成记录
- 查看历史详情
- 支持版本回滚

## API 接口

### 数据库相关

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/db/test` | POST | 测试数据库连接 |
| `/api/db/tables` | POST | 获取表列表 |
| `/api/db/table/detail` | POST | 获取表详情 |

### 代码生成

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/generate` | POST | 生成代码 |
| `/api/generate/preview` | POST | 预览代码 |
| `/api/generate/download` | GET | 下载代码 |

### 历史记录

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/history` | GET | 获取历史列表 |
| `/api/history/:id` | GET | 获取历史详情 |
| `/api/history/rollback` | POST | 回滚 |
| `/api/history/:id` | DELETE | 删除历史 |

## 请求示例

### 测试数据库连接

```bash
curl -X POST http://localhost:8199/api/db/test \
  -H "Content-Type: application/json" \
  -d '{"dsn":"mysql:root:123456@tcp(127.0.0.1:3306)/gfrd","type":"mysql"}'
```

### 获取表列表

```bash
curl -X POST http://localhost:8199/api/db/tables \
  -H "Content-Type: application/json" \
  -d '{"dsn":"mysql:root:123456@tcp(127.0.0.1:3306)/gfrd","type":"mysql"}'
```

### 生成代码

```bash
curl -X POST http://localhost:8199/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "dsn": "mysql:root:123456@tcp(127.0.0.1:3306)/gfrd",
    "type": "mysql",
    "tables": ["sys_dept", "sys_user"],
    "module": "sys",
    "output": "./server",
    "web": "./web",
    "features": ["list", "add", "edit", "delete", "view"]
  }'
```

## 与 CLI 的关系

Web 后台与 CLI 共用同一套核心代码：

- **parser** - 数据库解析器
- **generator** - 代码生成器
- **history** - 历史记录管理
- **template** - 代码模板

所有通过 Web 生成的代码与 CLI 生成的代码完全一致。

## 技术栈

- **后端**: GoFrame 2.7.1+
- **前端**: Vue 3 + Element Plus
- **数据库**: MySQL / PostgreSQL

## 目录结构

```
gen/
├── web/
│   ├── server.go        # Web 服务器
│   ├── handler.go       # API 处理器
│   └── static/
│       └── index.html   # 前端页面
├── cmd/                 # CLI 命令
├── generator/           # 生成器核心
├── parser/              # 解析器
├── history/             # 历史记录
└── template/            # 代码模板
```

## 注意事项

1. 首次启动前确保数据库连接正常
2. 生成代码前请先测试数据库连接
3. 回滚操作会覆盖当前文件，请先提交到 git
4. 默认端口为 8199，如需修改请编辑 `web/server.go`

## 常见问题

**Q: 启动后无法访问？**
A: 检查防火墙设置，确保 8199 端口开放。

**Q: 连接数据库失败？**
A: 检查 DSN 格式是否正确，数据库服务是否运行。

**Q: 生成的代码在哪里？**
A: 默认在 `./server` 和 `./web` 目录，可在生成配置中修改。

---

**版本**: 1.0.0
**日期**: 2026-02-27
