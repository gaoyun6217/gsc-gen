# gsc-gen

代码生成器项目集合，包含多个快速开发框架和模板。

## 项目结构

| 项目 | 描述 | 状态 |
|------|------|------|
| [gfrd/](gfrd/) | GFRD - GoFrame Rapid Development Framework | ✅ 已完成 |
| [hotgo-2.0/](hotgo-2.0/) | HotGo - 全栈管理系统 | ✅ 已完成 |
| [soybean-admin/](soybean-admin/) | SoybeanAdmin - Vue3 管理模板 | ✅ 已完成 |
| [gf2/](gf2/) | GoFrame 2 基础模板 | ✅ 已完成 |

---

## GFRD - GoFrame 快速开发框架

GFRD 是一个基于 GoFrame 2 的代码生成器，可以根据数据库表结构自动生成完整的 CRUD 代码。

### 特性

- **智能字段推断**: 自动识别字段类型并推断表单组件类型
- **完整 CRUD 生成**: 生成前后端完整的增删改查代码
- **树形结构支持**: 自动识别 `parent_id` 字段支持树形表格
- **软删除支持**: 自动识别 `deleted_at` 字段实现软删除
- **模板引擎**: 基于 text/template 的可扩展模板系统

### 技术栈

**后端**:
- GoFrame 2.7.1+
- Go 1.18+
- MySQL / PostgreSQL

**前端**:
- Vue 3 + TypeScript
- NaiveUI 组件库
- Vite 构建工具

### 目录结构

```
gfrd/
├── cmd/                  # CLI 命令
├── config/               # 配置管理
├── engine/               # 模板渲染引擎
├── generator/            # 代码生成器核心
├── parser/               # 数据库表解析
├── types/                # 类型定义
├── template/             # 代码模板
│   ├── backend/          # 后端模板 (Go)
│   ├── frontend/         # 前端模板 (Vue/TS)
│   └── sql/              # SQL 模板
├── server/               # 生成的后端项目
└── web/                  # 生成的前端项目
```

### 快速开始

1. **配置数据库连接**

编辑 `gfrd/config/config.yaml`:

```yaml
database:
  default:
    link: "mysql:username:password@tcp(127.0.0.1:3306)/database"
```

2. **运行代码生成器**

```bash
cd gfrd
go run cmd/cmd.go
```

3. **查看生成的代码**

- 后端代码：`gfrd/server/`
- 前端代码：`gfrd/web/`

### 生成的文件

**后端 (4 个文件)**:
- `api/sys/sys_dept.go` - API 定义
- `internal/handler/sys/sys_dept.go` - Handler 实现
- `internal/router/genrouter/sys_dept.go` - 路由注册
- `storage/data/generate/sys_dept_menu.sql` - 菜单 SQL

**前端 (4 个文件)**:
- `src/api/sys/dept/index.ts` - API 服务
- `src/api/sys/dept/types.ts` - TypeScript 类型
- `src/views/sys/dept/index.vue` - 列表页面
- `src/views/sys/dept/edit.vue` - 编辑弹窗

### 字段类型推断规则

| 字段名/类型 | 推断的表单类型 |
|-------------|----------------|
| `status` | switch (开关) |
| `remark`, `description` | textarea (文本域) |
| `email` | input (邮箱) |
| `phone`, `mobile` | input (电话) |
| `password` | input (密码) |
| `date`, `time` | input (日期时间) |
| `tinyint(1)` | switch (布尔值) |

### 查询条件推断

以下字段名会自动添加为查询条件：
- `name` - 名称
- `code` - 编码
- `status` - 状态
- `type` - 类型
- `keyword` - 关键词

### 文档

- [使用指南](gfrd/GUIDE.md) - 详细使用说明
- [测试报告](gfrd/TEST_REPORT.md) - 代码生成测试报告

---

## HotGo 2.0

完整的全栈管理系统，基于 GoFrame 2.9.4 和 Vue3 + NaiveUI。

### 技术栈

- **后端**: GoFrame 2.9.4, Go 1.24+, Casbin, JWT
- **前端**: Vue3, TypeScript, NaiveUI, Pinia

### 特性

- 多应用入口（后台管理、前台、API、WebSocket）
- 插件化系统
- RBAC 权限管理
- 内置代码生成器

[查看 HotGo 详细文档](hotgo-2.0/README.md)

---

## SoybeanAdmin

现代化的 Vue3 管理后台模板。

### 技术栈

- Vue 3.5+
- Vite 7+
- TypeScript 5+
- NaiveUI
- UnoCSS
- Pinia
- Vue Router 4

### 特性

- 文件路由（Elegant Router）
- 主题系统
- 动态路由权限
- Mock 数据支持

[查看 SoybeanAdmin 详细文档](soybean-admin/README.md)

---

## gf2

GoFrame 2 基础项目模板。

### 技术栈

- GoFrame 2.7.1
- Go 1.18+

---

## 开发工具

### 依赖安装

```bash
# 安装 GoFrame CLI 工具
go install github.com/gogf/gf/cmd/gf/v2@latest

# 安装 Node.js 依赖
cd hotgo-2.0/web && pnpm install
cd soybean-admin && pnpm install
```

### 常用命令

**GFRD 生成器**:
```bash
cd gfrd
go run cmd/cmd.go
```

**HotGo 后端**:
```bash
cd hotgo-2.0/server
go run main.go
```

**HotGo 前端**:
```bash
cd hotgo-2.0/web
pnpm dev
```

**SoybeanAdmin**:
```bash
cd soybean-admin
pnpm dev
```

---

## License

MIT License

---

**作者**: AI Assistant
**创建日期**: 2026-02-27
