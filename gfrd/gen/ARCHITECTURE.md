# GFRD Gen 模块架构设计文档

## 1. 模块概述

GFRD Gen 是一个独立的代码生成器模块，用于为 GoFrame 2 + SoybeanAdmin 框架生成全栈 CRUD 代码。

### 1.1 设计目标

- **独立性**: 可作为独立 package 维护和升级
- **可扩展性**: 支持自定义模板和生成逻辑
- **易用性**: 简单的 CLI 命令即可生成完整代码
- **智能化**: 自动推断字段类型、表单类型、验证规则

## 2. 目录结构

```
gen/
├── main.go                    # CLI 入口文件
├── go.mod                     # 模块依赖
├── Makefile                   # 构建脚本
├── README.md                  # 使用说明
├── .gitignore                 # Git 忽略文件
├── generator.yaml.example     # 配置示例
│
├── cmd/                       # CLI 命令层
│   └── cmd.go                 # 命令定义和执行
│
├── config/                    # 配置管理
│   └── config.go              # 配置文件加载/保存
│
├── types/                     # 类型定义
│   ├── types.go               # 核心类型定义
│   └── namecase.go            # 名称转换工具
│
├── parser/                    # 数据库解析器
│   └── parser.go              # 表结构解析
│
├── engine/                    # 模板渲染引擎
│   └── renderer.go            # 模板渲染和文件输出
│
├── generator/                 # 生成器核心
│   └── generator.go           # 生成逻辑编排
│
└── template/                  # 模板文件
    ├── backend/
    │   ├── api.go.tpl         # API 定义模板
    │   ├── handler.go.tpl     # Handler 实现模板
    │   ├── service.go.tpl     # Service 接口模板
    │   ├── router.go.tpl      # 路由注册模板
    │   └── test.go.tpl        # 单元测试模板
    │
    ├── frontend/
    │   ├── api.ts.tpl         # API 服务模板
    │   ├── types.ts.tpl       # TS 类型定义模板
    │   ├── index.vue.tpl      # 列表页模板
    │   └── edit.vue.tpl       # 编辑弹窗模板
    │
    └── sql/
        └── menu.sql.tpl       # 菜单 SQL 模板
```

## 3. 核心组件

### 3.1 CLI 命令层 (cmd/)

负责解析命令行参数并调用生成器。

**主要命令**:
- `crud` - 生成完整 CRUD 代码
- `backend` - 仅生成后端代码
- `frontend` - 仅生成前端代码
- `preview` - 预览生成结果

**命令参数**:
```go
type Config struct {
    Table        string   // 表名
    DB           string   // 数据库连接
    Output       string   // 后端输出目录
    WebOutput    string   // 前端输出目录
    Package      string   // Go 包名
    Module       string   // 模块名
    Features     []string // 要生成的功能
    WithTest     bool     // 是否生成测试
    WithDoc      bool     // 是否生成文档
    LayerMode    string   // 分层模式
    Preview      bool     // 是否仅预览
}
```

### 3.2 类型定义层 (types/)

定义生成器使用的核心数据结构。

**TableInfo** - 表结构信息:
```go
type TableInfo struct {
    Name        string        // 表名
    Comment     string        // 表注释
    Columns     []*ColumnInfo // 列信息
    PrimaryKey  string        // 主键列名
    Indexes     []*IndexInfo  // 索引信息
    IsTreeTable bool          // 是否为树形表
}
```

**ColumnInfo** - 列信息:
```go
type ColumnInfo struct {
    Name         string // 列名 (下划线)
    NameCamel    string // 列名 (小驼峰)
    NamePascal   string // 列名 (大驼峰)
    Type         string // 数据库类型
    TypeGo       string // Go 类型
    TypeTs       string // TypeScript 类型
    Comment      string // 注释
    Nullable     bool   // 是否可空
    IsPrimary    bool   // 是否主键
    IsListField  bool   // 是否在列表中显示
    IsQueryField bool   // 是否作为查询条件
    FormType     string // 表单类型
    QueryType    string // 查询类型
}
```

### 3.3 数据库解析器 (parser/)

负责连接数据库并解析表结构。

**核心方法**:
- `New(dsn, driver)` - 创建解析器
- `ParseTable(ctx, tableName)` - 解析表结构
- `parseTableMySQL(ctx, tableName)` - MySQL 解析
- `parseTablePostgres(ctx, tableName)` - PostgreSQL 解析
- `dataTypeToGo(dataType)` - 数据库类型转 Go 类型
- `dataTypeToTs(dataType)` - 数据库类型转 TypeScript 类型
- `inferFormType(dataType, comment)` - 推断表单类型

**类型映射**:

| 数据库类型 | Go 类型 | TypeScript 类型 |
|-----------|---------|----------------|
| int/bigint | int64 | number |
| int | int | number |
| varchar/text | string | string |
| datetime/timestamp | *gtime.Time | string |
| bool/tinyint(1) | bool | boolean |
| decimal/numeric | string | number |
| json | gjson.RawMessage | any |

### 3.4 模板渲染引擎 (engine/)

负责加载模板文件并渲染生成代码。

**核心方法**:
- `NewRenderer(templateDir)` - 创建渲染器
- `Render(ctx, tmplFile, data)` - 渲染模板
- `RenderAndWrite(ctx, tmplFile, outputPath, data)` - 渲染并写入文件

**模板函数**:
```go
funcMap = template.FuncMap{
    // 名称转换
    "toCamel":  types.ToCamel,
    "toPascal": types.ToPascal,
    "toKebab":  types.ToKebab,
    "toSnake":  types.ToSnake,

    // 字符串操作
    "lower":    strings.ToLower,
    "upper":    strings.ToUpper,
    "contains": strings.Contains,
    "replace":  strings.ReplaceAll,

    // 字段过滤
    "filterQueryFields": filterQueryFields,
    "filterListFields":  filterListFields,
    "filterFormFields":  filterFormFields,
}
```

### 3.5 生成器核心 (generator/)

负责编排整个生成流程。

**核心流程**:
```go
func (g *Generator) Generate(ctx context.Context) error {
    // 1. 初始化数据库解析器
    p, err := parser.New(dsn, dbType)

    // 2. 解析表结构
    table, err := p.ParseTable(ctx, g.cfg.Table)

    // 3. 准备渲染数据
    renderData := g.prepareRenderData(table)

    // 4. 生成后端代码
    g.generateBackend(ctx, renderData)

    // 5. 生成前端代码
    g.generateFrontend(ctx, renderData)

    return nil
}
```

## 4. 模板系统

### 4.1 模板变量

所有模板都可以访问以下变量:

```go
type RenderData struct {
    Table          *TableInfo      // 表信息
    Package        string          // 包名
    Module         string          // 模块名
    EntityName     string          // 实体名 (Pascal)
    EntityKebab    string          // 实体名 (kebab-case)
    EntitySnake    string          // 实体名 (snake_case)
    Operations     []*OperationInfo // 操作列表
    Features       map[string]bool // 功能开关
    HasTree        bool            // 是否有树结构
    HasSoftDelete  bool            // 是否有软删除
    HasCreatedAt   bool            // 是否有创建时间
    HasUpdatedAt   bool            // 是否有更新时间
}
```

### 4.2 后端模板

**api.go.tpl** - 生成 API 定义:
- ListReq/ListRes - 列表请求/响应
- ViewReq/ViewRes - 详情请求/响应
- AddReq - 新增请求
- EditReq - 修改请求
- DeleteReq - 删除请求

**handler.go.tpl** - 生成 Handler 实现:
- List() - 列表查询
- View() - 详情查看
- Add() - 新增
- Edit() - 修改
- Delete() - 删除

**service.go.tpl** - 生成 Service 接口 (standard 模式):
- I{Entity}Service 接口定义
- {Entity}Service 实现

**router.go.tpl** - 生成路由注册:
- 路由分组
- 路由绑定

### 4.3 前端模板

**api.ts.ts** - 生成 API 服务:
- {EntityName}List() - 列表 API
- {EntityName}View() - 详情 API
- {EntityName}Add() - 新增 API
- {EntityName}Edit() - 修改 API
- {EntityName}Delete() - 删除 API

**types.ts.tpl** - 生成 TypeScript 类型:
- {EntityName} 接口
- {EntityName}EditDTO 接口
- {EntityName}ListParams 接口

**index.vue.tpl** - 生成列表页:
- 搜索表单
- 数据表格
- 操作列
- 分页

**edit.vue.tpl** - 生成编辑弹窗:
- 表单
- 验证规则
- 提交逻辑

## 5. 使用示例

### 5.1 基本使用

```bash
# 安装
go install github.com/gfrd/gen@latest

# 生成 CRUD
gfrd-gen crud \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web" \
  --module="sys"
```

### 5.2 使用配置文件

```bash
gfrd-gen crud --config=generator.yaml
```

### 5.3 预览模式

```bash
gfrd-gen preview \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd"
```

## 6. 扩展开发

### 6.1 添加新模板

1. 在 `template/backend/` 或 `template/frontend/` 创建新模板文件
2. 在 `generator/generator.go` 中添加生成逻辑

### 6.2 自定义类型映射

修改 `parser/parser.go` 中的 `dataTypeToGo` 和 `dataTypeToTs` 方法。

### 6.3 自定义表单推断

修改 `parser/parser.go` 中的 `inferFormType` 方法。

## 7. 最佳实践

### 7.1 模板开发

- 使用 `{{- range }}` 而不是 `{{ range }}` 来避免多余换行
- 使用 `{{- if }}` 来控制空白
- 在模板开头添加 `// Code generated by gfrd-gen. DO NOT EDIT.` 注释

### 7.2 字段配置

- 在数据库表注释中使用"状态"、"是否"等词会自动推断为 switch 表单类型
- 使用"性别"、"类型"等词会自动推断为 radio 表单类型
- 使用"图片"、"头像"等词会自动推断为 upload 表单类型

### 7.3 查询字段

以下字段名会自动成为查询条件:
- name, username, code - 名称类
- status, type - 状态类型类
- email, phone, mobile - 联系方式类

## 8. 依赖管理

### 8.1 Go 依赖

```go
require (
    github.com/go-sql-driver/mysql v1.7.1
    github.com/gogf/gf/v2 v2.7.1
    github.com/spf13/cobra v1.8.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### 8.2 版本兼容性

- Go 1.21+
- GoFrame 2.7+
- Node.js 18+ (前端生成代码)

## 9. 测试

```bash
# 运行测试
make test

# 生成示例代码
make demo

# 预览生成结果
make preview
```

## 10. 未来计划

- [ ] 支持更多数据库类型 (Oracle, SQL Server)
- [ ] 支持主子表生成
- [ ] 支持导入/导出功能
- [ ] 支持 GraphQL API
- [ ] Web UI 配置界面
- [ ] AI 辅助字段推断
