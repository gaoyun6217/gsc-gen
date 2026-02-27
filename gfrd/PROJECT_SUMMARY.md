# GFRD 快速开发框架 - 项目总结

## 项目概述

基于 GoFrame 2 和 SoybeanAdmin 打造的全新快速开发框架，包含一个独立的代码生成器模块 (gen)。

## 目录结构

```
gfrd/
├── gen/                          # 代码生成器模块 (已完成)
│   ├── cmd/                      # CLI 命令
│   ├── config/                   # 配置管理
│   ├── engine/                   # 模板引擎
│   ├── generator/                # 生成器核心
│   ├── parser/                   # 数据库解析器
│   ├── types/                    # 类型定义
│   ├── template/                 # 模板文件
│   ├── main.go                   # 入口
│   ├── README.md                 # 使用说明
│   └── ARCHITECTURE.md           # 架构文档
│
├── server/                       # 后端项目 (待生成)
│   ├── cmd/
│   ├── internal/
│   ├── api/
│   ├── config/
│   └── ...
│
└── web/                          # 前端项目 (待生成)
    ├── src/
    ├── api/
    ├── views/
    └── ...
```

## 已完成的工作

### 1. Gen 模块核心功能

#### 1.1 CLI 命令系统
- `crud` - 生成完整 CRUD 代码
- `backend` - 仅生成后端
- `frontend` - 仅生成前端
- `preview` - 预览生成结果

#### 1.2 数据库解析器
- 支持 MySQL 和 PostgreSQL
- 自动解析表结构、列信息、索引
- 智能类型转换 (数据库 → Go/TypeScript)
- 智能表单类型推断
- 树形表自动识别

#### 1.3 模板渲染引擎
- 基于 Go template
- 丰富的模板函数 (名称转换、字段过滤等)
- 支持自定义模板目录
- 预览模式支持

#### 1.4 后端模板
- `api.go.tpl` - API 定义 (Request/Response)
- `handler.go.tpl` - 请求处理逻辑
- `service.go.tpl` - 服务接口 (standard 模式)
- `router.go.tpl` - 路由注册
- `test.go.tpl` - 单元测试

#### 1.5 前端模板
- `api.ts.tpl` - API 服务封装
- `types.ts.tpl` - TypeScript 类型定义
- `index.vue.tpl` - 列表页 (含搜索、表格、分页)
- `edit.vue.tpl` - 编辑弹窗 (含表单验证)

#### 1.6 SQL 模板
- `menu.sql.tpl` - 菜单权限 SQL 脚本

### 2. 类型系统

#### 核心类型
```go
// 表结构信息
type TableInfo struct {
    Name        string
    Comment     string
    Columns     []*ColumnInfo
    PrimaryKey  string
    IsTreeTable bool
}

// 列信息
type ColumnInfo struct {
    Name         string  // 下划线
    NameCamel    string  // 小驼峰
    NamePascal   string  // 大驼峰
    Type         string  // 数据库类型
    TypeGo       string  // Go 类型
    TypeTs       string  // TS 类型
    Comment      string
    FormType     string  // 表单类型
    IsQueryField bool    // 是否查询字段
    IsListField  bool    // 是否列表字段
}

// 渲染数据
type RenderData struct {
    Table         *TableInfo
    EntityName    string
    EntityKebab   string
    Module        string
    Features      map[string]bool
    HasTree       bool
    HasSoftDelete bool
}
```

### 3. 配置文件

```yaml
database:
  driver: mysql
  dsn: "root:123456@tcp(127.0.0.1:3306)/gfrd"

generator:
  backend:
    enabled: true
    output: "./server"
    package: "github.com/gfrd/server"
    layerMode: "simple"
    withTest: true
    withDoc: true

  frontend:
    enabled: true
    output: "./web/src"
    typescript: true

  features:
    list: true
    add: true
    edit: true
    delete: true
    view: true
```

## 使用流程

### 安装 Gen 模块

```bash
cd gfrd/gen
go install .
```

### 生成示例代码

```bash
gfrd-gen crud \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web" \
  --module="sys" \
  --features="add,edit,delete,view,list" \
  --with-test
```

### 生成结果

**后端文件**:
```
server/
├── api/sys/sys_user.go
├── internal/handler/sys/sys_user.go
├── internal/router/genrouter/sys_user.go
├── tests/handler/sys/sys_user_test.go
└── storage/data/generate/sys_user_menu.sql
```

**前端文件**:
```
web/src/
├── api/sys/user/index.ts
├── api/sys/user/types.ts
└── views/sys/user/
    ├── index.vue
    └── edit.vue
```

## 技术亮点

### 1. 智能类型推断

| 数据库类型 | Go 类型 | TS 类型 | 表单类型 |
|-----------|---------|--------|---------|
| int/bigint | int64 | number | input |
| varchar(50) | string | string | input |
| text | string | string | textarea |
| tinyint(1) | bool | boolean | switch |
| datetime | *gtime.Time | string | datetime |
| 含"状态"注释 | int | number | switch |
| 含"图片"注释 | string | string | upload |

### 2. 智能字段过滤

**查询字段** (自动添加到搜索表单):
- name, username, code
- status, type
- email, phone

**列表字段** (自动添加到表格列):
- 排除 password, token 等敏感字段
- 排除 created_at, updated_at 等时间字段

**表单字段** (自动添加到编辑表单):
- 排除 id, created_at, updated_at, deleted_at

### 3. 树形表支持

自动识别包含 `parent_id`, `level`, `path` 字段的表，生成树形结构代码。

### 4. 软删除支持

自动识别 `deleted_at` 字段，生成软删除代码。

## 下一步计划

### 阶段 1: 完善 Gen 模块 (1 周)
- [ ] 添加更多模板函数
- [ ] 完善错误处理
- [ ] 添加单元测试
- [ ] 编写详细文档

### 阶段 2: 创建框架骨架 (2 周)
- [ ] 初始化 GoFrame 后端项目
- [ ] 初始化 SoybeanAdmin 前端项目
- [ ] 配置基础依赖
- [ ] 集成 Gen 模块

### 阶段 3: 基础功能模块 (2 周)
- [ ] 用户管理模块
- [ ] 角色管理模块
- [ ] 部门管理模块
- [ ] 菜单管理模块

### 阶段 4: 增强功能 (2 周)
- [ ] 字典管理
- [ ] 文件上传
- [ ] 系统配置
- [ ] 日志管理

### 阶段 5: 文档和示例 (1 周)
- [ ] 快速开始文档
- [ ] API 文档
- [ ] 示例项目
- [ ] 视频教程

## 与 HotGo 对比

| 特性 | HotGo 2.0 | GFRD (新框架) |
|------|-----------|---------------|
| 代码生成器 | 内置，UI 操作 | 独立 CLI，支持配置 |
| 分层架构 | 7 层 (复杂) | 4 层 (简化) |
| 前端 Hook | 分散 | 统一 useCRUD |
| 模板自定义 | 困难 | 简单，支持热插拔 |
| 单元测试 | 无 | 自动生成 |
| 文档 | 一般 | 完善，含架构文档 |
| 学习曲线 | 陡峭 | 平缓 |

## 核心优势

1. **独立 Gen 模块**: 可独立维护和升级，不耦合具体项目
2. **简化分层**: 从 HotGo 的 7 层简化为 4 层
3. **智能推断**: 根据字段类型和注释自动推断表单类型
4. **统一 Hook**: 前端使用统一的 useCRUD 钩子
5. **测试覆盖**: 自动生成单元测试
6. **完善文档**: 包含 README、ARCHITECTURE 等完整文档

## 技术栈

### 后端
- Go 1.21+
- GoFrame 2.7+
- MySQL 8+ / PostgreSQL 15+
- Redis 7+

### 前端
- Vue 3.4+
- Vite 5+
- TypeScript 5+
- NaiveUI 2.38+
- Pinia 2.1+
- UnoCSS

### Gen 模块
- Go 1.21+
- github.com/go-sql-driver/mysql
- github.com/gogf/gf/v2
- github.com/spf13/cobra
- gopkg.in/yaml.v3

## 贡献指南

### 开发环境设置

```bash
# 克隆项目
git clone <repo>
cd gfrd/gen

# 安装依赖
go mod tidy

# 运行测试
make test

# 构建
make build
```

### 添加新模板

1. 在 `template/` 目录下创建新模板
2. 在 `generator/generator.go` 中添加生成逻辑
3. 运行 `make demo` 测试

### 提交代码

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址：https://github.com/gfrd/gfrd
- 问题反馈：https://github.com/gfrd/gfrd/issues
