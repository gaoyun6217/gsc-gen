# GFRD 框架创建总结

## 创建时间
2026-02-27

## 项目概述

成功创建了一个基于 GoFrame 2 和 SoybeanAdmin 的全新快速开发框架，包含一个独立的代码生成器模块 (gen)。

## 目录结构

```
gfrd/
├── go.work                     # Go Workspace 配置
│
├── gen/                        # 代码生成器模块（独立可维护）
│   ├── main.go                 # CLI 入口
│   ├── go.mod                  # 模块依赖
│   ├── Makefile                # 构建脚本
│   ├── README.md               # 使用说明
│   ├── QUICKSTART.md           # 快速开始指南
│   ├── ARCHITECTURE.md         # 架构设计文档
│   │
│   ├── cmd/                    # CLI 命令层
│   │   └── cmd.go              # crud/backend/frontend/preview 命令
│   │
│   ├── config/                 # 配置管理
│   │   └── config.go           # YAML 配置加载/保存
│   │
│   ├── types/                  # 类型定义
│   │   ├── types.go            # TableInfo, ColumnInfo 等
│   │   └── namecase.go         # 名称转换工具
│   │
│   ├── parser/                 # 数据库解析器
│   │   └── parser.go           # MySQL/PostgreSQL 解析
│   │
│   ├── engine/                 # 模板渲染引擎
│   │   └── renderer.go         # Template 渲染
│   │
│   ├── generator/              # 生成器核心
│   │   └── generator.go        # 生成逻辑编排
│   │
│   └── template/               # 模板文件
│       ├── backend/
│       │   ├── api.go.tpl      # API 定义
│       │   ├── handler.go.tpl  # Handler 实现
│       │   ├── service.go.tpl  # Service 接口
│       │   ├── router.go.tpl   # 路由注册
│       │   └── test.go.tpl     # 单元测试
│       │
│       ├── frontend/
│       │   ├── api.ts.tpl      # API 服务
│       │   ├── types.ts.tpl    # TS 类型
│       │   ├── index.vue.tpl   # 列表页
│       │   └── edit.vue.tpl    # 编辑弹窗
│       │
│       └── sql/
│           └── menu.sql.tpl    # 菜单 SQL
│
├── server/                     # GoFrame 后端项目
│   ├── cmd/
│   │   └── main.go             # 入口文件
│   ├── config/
│   │   └── config.yaml         # 配置文件
│   ├── go.mod                  # 依赖管理
│   ├── api/                    # API 定义（生成）
│   ├── internal/
│   │   ├── handler/            # 处理层（生成）
│   │   ├── service/            # 服务层（生成）
│   │   ├── model/
│   │   ├── router/             # 路由（生成）
│   │   └── middleware/         # 中间件
│   └── storage/                # 运行时数据
│
├── web/                        # SoybeanAdmin 前端项目
│   ├── src/
│   │   ├── api/                # API 服务（生成）
│   │   ├── views/              # 页面组件（生成）
│   │   ├── components/         # 公共组件
│   │   ├── composables/        # 组合式函数
│   │   │   └── useCRUD.ts      # CRUD 钩子
│   │   ├── router/             # 路由配置
│   │   ├── stores/             # Pinia 状态
│   │   └── utils/              # 工具函数
│   ├── package.json            # 依赖管理
│   ├── vite.config.js          # Vite 配置
│   └── index.html              # 入口文件
│
├── README.md                   # 项目说明
├── GUIDE.md                    # 使用指南
├── PROJECT_SUMMARY.md          # 项目总结
├── test.bat                    # Windows 测试脚本
└── test.sh                     # Linux/Mac 测试脚本
```

## 核心功能

### 1. Gen 模块（代码生成器）

**CLI 命令**:
- `crud` - 生成完整 CRUD 代码
- `backend` - 仅生成后端
- `frontend` - 仅生成前端
- `preview` - 预览生成结果

**支持的数据库**:
- MySQL 8+
- PostgreSQL 15+

**智能功能**:
- 自动类型转换（数据库 → Go/TypeScript）
- 智能表单类型推断
- 自动识别查询字段
- 树形表支持
- 软删除支持

**生成的代码**:
- 后端：API、Handler、Service、Router、Test
- 前端：API 服务、TS 类型、Vue 组件
- SQL：菜单权限脚本

### 2. Server 项目（GoFrame 后端）

**技术栈**:
- Go 1.23+
- GoFrame 2.10+
- MySQL/PostgreSQL
- Redis

**目录约定**:
- `api/` - API 输入/输出定义
- `internal/handler/` - 请求处理
- `internal/service/` - 业务逻辑
- `internal/router/` - 路由注册
- `config/` - 配置文件

### 3. Web 项目（SoybeanAdmin 前端）

**技术栈**:
- Vue 3.4+
- Vite 5+
- TypeScript 5+
- NaiveUI 2.38+
- Pinia 2.1+

**核心功能**:
- 统一的 request 工具（Axios 封装）
- useCRUD 组合式函数
- 路由 guards
- 权限管理

## 使用方法

### 从根目录运行（推荐）

```bash
# 生成 CRUD 代码
go run gen/main.go crud \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web/src" \
  --module="sys" \
  --features="add,edit,delete,view,list" \
  --with-test
```

### 安装为全局命令

```bash
cd gen
go install .
gfrd-gen crud --table="sys_user" --db="mysql:..."
```

### 在代码中导入

```go
import "github.com/gfrd/gen/generator"

cfg := &generator.Config{
    Table:  "sys_user",
    DB:     "mysql:...",
    Output: "./server",
}
generator.NewGenerator(cfg).Generate(ctx)
```

## 测试状态

✅ 所有测试通过

```
[OK] gen 模块存在
[OK] gen 模块构建成功
[OK] gen CLI 执行成功
[OK] server 目录存在
[OK] web 目录存在
[OK] 从根目录运行 gen 成功
```

## 技术亮点

1. **Go Workspace**: 使用 go.work 管理多模块项目
2. **独立 Gen 模块**: 可作为独立 package 维护和升级
3. **智能生成**: 根据表结构自动推断代码逻辑
4. **全栈生成**: 同时生成后端和前端代码
5. **可定制模板**: 支持自定义模板文件

## 下一步计划

1. 完善基础模块（用户、角色、部门、菜单管理）
2. 添加更多前端 UI 组件
3. 集成 Casbin 权限管理
4. 添加 JWT 认证
5. 完善文档和示例

## 文件统计

| 类型 | 数量 |
|------|------|
| Go 文件 | 10 |
| TypeScript 文件 | 3 |
| Vue 文件 | 1 |
| 模板文件 | 9 |
| 配置文件 | 5 |
| 文档文件 | 6 |

## 依赖版本

| 模块 | 版本 |
|------|------|
| Go | 1.23.0 |
| GoFrame | 2.10.0 |
| Vue | 3.4.0 |
| TypeScript | 5.3.0 |
| NaiveUI | 2.38.0 |

---

**项目已成功创建并可投入使用！**
