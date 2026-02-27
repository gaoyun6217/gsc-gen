# GFRD 快速开发框架

基于 GoFrame 2 和 SoybeanAdmin 的全新快速开发框架。

## 目录结构

```
gfrd/
├── gen/                    # 代码生成器模块
├── server/                 # GoFrame 后端
├── web/                    # SoybeanAdmin 前端
├── go.mod                  # 根模块依赖
└── README.md               # 项目说明
```

## 快速开始

### 1. 初始化项目

```bash
# 安装后端依赖
cd server
go mod tidy

# 安装前端依赖
cd web
pnpm install
```

### 2. 使用代码生成器

```bash
# 从根目录运行
cd gfrd

# 生成用户管理模块
go run ./gen crud \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web/src" \
  --module="sys"
```

### 3. 启动服务

```bash
# 后端
cd server
go run main.go

# 前端
cd web
pnpm dev
```

## 使用 gen 模块

### 作为 CLI 工具

```bash
# 方式 1: 直接运行
go run ./gen crud --table="sys_user" --db="mysql:..."

# 方式 2: 安装后使用
cd gen
go install .
gfrd-gen crud --table="sys_user" --db="mysql:..."
```

### 作为代码库导入

```go
package main

import (
    "context"
    "github.com/gfrd/gen/generator"
)

func main() {
    cfg := &generator.Config{
        Table:   "sys_user",
        DB:      "mysql:...",
        Output:  "./server",
        Module:  "sys",
    }

    gen := generator.NewGenerator(cfg)
    gen.Generate(context.Background())
}
```

## 技术栈

### 后端
- Go 1.21+
- GoFrame 2.6+
- MySQL 8+ / PostgreSQL 15+
- Redis 7+

### 前端
- Vue 3.4+
- Vite 5+
- TypeScript 5+
- NaiveUI 2.38+
- Pinia 2.1+

## 许可证

MIT License
