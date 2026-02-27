# GFRD 框架使用指南

## 项目结构

```
gfrd/
├── go.work                 # Go Workspace 配置
├── gen/                    # 代码生成器模块
├── server/                 # GoFrame 后端
├── web/                    # SoybeanAdmin 前端
├── test.bat                # Windows 测试脚本
└── test.sh                 # Linux/Mac 测试脚本
```

## 环境要求

- Go 1.23+
- Node.js 18+
- MySQL 8+ / PostgreSQL 15+
- pnpm 10+

## 快速开始

### 1. 安装依赖

```bash
# 后端
cd server
go mod tidy

# 前端
cd web
pnpm install
```

### 2. 配置数据库

创建数据库：

```sql
CREATE DATABASE gfrd DEFAULT CHARACTER SET utf8mb4;
```

### 3. 使用代码生成器

从项目根目录运行：

```bash
# 预览生成结果
go run gen/main.go preview \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --module="sys"

# 生成完整 CRUD
go run gen/main.go crud \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web/src" \
  --module="sys" \
  --features="add,edit,delete,view,list" \
  --with-test
```

### 4. 启动服务

```bash
# 后端
cd server
go run cmd/main.go

# 前端
cd web
pnpm dev
```

## 使用 gen 模块

### 方式 1: 从根目录运行（推荐）

```bash
go run gen/main.go crud --table="sys_user" --db="mysql:..."
```

### 方式 2: 安装为全局命令

```bash
cd gen
go install .
gfrd-gen crud --table="sys_user" --db="mysql:..."
```

### 方式 3: 在代码中导入

```go
package main

import (
    "context"
    "github.com/gfrd/gen/generator"
)

func main() {
    cfg := &generator.Config{
        Table:   "sys_user",
        DB:      "mysql:root:123456@tcp(127.0.0.1:3306)/gfrd",
        Output:  "./server",
        Module:  "sys",
        Features: []string{"add", "edit", "delete", "view", "list"},
    }

    gen := generator.NewGenerator(cfg)
    gen.Generate(context.Background())
}
```

## 运行测试

### Windows

```bash
test.bat
```

### Linux/Mac

```bash
chmod +x test.sh
./test.sh
```

## 生成的代码结构

生成 `sys_user` 表后：

```
server/
├── api/sys/
│   └── sys_user.go              # API 定义
├── internal/handler/sys/
│   └── sys_user.go              # Handler 实现
├── internal/router/genrouter/
│   └── sys_user.go              # 路由注册
└── storage/data/generate/
    └── sys_user_menu.sql        # 菜单 SQL

web/src/
├── api/sys/user/
│   ├── index.ts                 # API 服务
│   └── types.ts                 # TS 类型
└── views/sys/user/
    ├── index.vue                # 列表页
    └── edit.vue                 # 编辑弹窗
```

## 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| --table, -t | 表名（必填） | - |
| --db, -d | 数据库 DSN（必填） | - |
| --output | 后端输出目录 | ./server |
| --web-output | 前端输出目录 | ./web |
| --package | Go 包名 | github.com/gfrd/server |
| --module, -m | 模块名 | sys |
| --features | 要生成的功能 | add,edit,delete,view,list |
| --with-test | 生成单元测试 | false |
| --with-doc | 生成 API 文档 | true |
| --layer-mode | 分层模式 | simple |
| --preview | 仅预览 | false |

## 常见问题

### Q: go.work 文件有什么用？

A: Go Workspace 允你在一个仓库中管理多个模块，方便本地开发和测试。

### Q: 如何添加新的生成模板？

A: 在 `gen/template/` 目录下添加新的模板文件，然后在 `gen/generator/generator.go` 中添加生成逻辑。

### Q: 如何修改类型映射？

A: 修改 `gen/parser/parser.go` 中的 `dataTypeToGo` 和 `dataTypeToTs` 方法。

## 许可证

MIT License
