# GFRD Code Generator

GoFrame Rapid Development Code Generator - 基于 GoFrame 2 和 SoybeanAdmin 的全栈代码生成器

## 功能特性

- ✅ 支持 MySQL / PostgreSQL 数据库
- ✅ 全栈代码生成（后端 API + Handler + 前端 Vue 组件）
- ✅ 智能字段类型推断（Go/TypeScript）
- ✅ 智能表单类型推断（input/textarea/select/date 等）
- ✅ 支持预览模式
- ✅ 支持单元测试生成
- ✅ 支持 API 文档注释
- ✅ 可配置的生成选项
- ✅ 支持树形表结构
- ✅ 支持软删除
- ✅ 菜单权限 SQL 生成

## 快速开始

### 1. 安装

```bash
cd gen
go install .
```

或者构建本地可执行文件：

```bash
make build
```

### 2. 使用

#### 生成完整 CRUD

```bash
gfrd-gen crud \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web" \
  --module="sys" \
  --features="add,edit,delete,view,list" \
  --with-test \
  --with-doc
```

#### 仅生成后端

```bash
gfrd-gen backend \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --module="sys"
```

#### 仅生成前端

```bash
gfrd-gen frontend \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --web-output="./web" \
  --module="sys"
```

#### 预览生成结果

```bash
gfrd-gen preview \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --module="sys"
```

### 3. 命令行参数

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| --table | -t | 表名（必填） | - |
| --db | -d | 数据库 DSN（必填） | - |
| --output | | 后端输出目录 | ./server |
| --web-output | | 前端输出目录 | ./web |
| --package | | Go 包名 | github.com/gfrd/server |
| --module | -m | 模块名 | sys |
| --features | | 要生成的功能 | add,edit,delete,view,list |
| --with-test | | 生成单元测试 | false |
| --with-doc | | 生成 API 文档 | true |
| --layer-mode | | 分层模式（simple/standard） | simple |
| --preview | | 仅预览，不写文件 | false |
| --template | | 自定义模板目录 | - |

### 4. 功能选项

`--features` 参数支持以下功能：

- `list` - 列表查询
- `add` - 新增
- `edit` - 修改
- `delete` - 删除
- `view` - 详情
- `export` - 导出
- `import` - 导入

## 配置文件

使用配置文件可以更方便地管理生成选项：

```bash
gfrd-gen crud --config=generator.yaml
```

配置文件示例（generator.yaml）：

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
    export: false
```

## 生成结果

### 后端文件结构

```
server/
├── api/sys/
│   └── sys_user.go          # API 定义
├── internal/handler/sys/
│   └── sys_user.go          # Handler 实现
├── internal/service/sys/
│   └── sys_user.go          # Service 接口（standard 模式）
└── tests/handler/sys/
    └── sys_user_test.go     # 单元测试
```

### 前端文件结构

```
web/src/
├── api/sys/user/
│   ├── index.ts             # API 服务
│   └── types.ts             # TypeScript 类型
└── views/sys/user/
    ├── index.vue            # 列表页
    └── edit.vue             # 编辑弹窗
```

## 自定义模板

在 `template/` 目录下可以自定义模板：

```
template/
├── backend/
│   ├── api.go.tpl
│   ├── handler.go.tpl
│   ├── service.go.tpl
│   └── test.go.tpl
├── frontend/
│   ├── api.ts.tpl
│   ├── types.ts.tpl
│   ├── index.vue.tpl
│   └── edit.vue.tpl
└── sql/
    └── menu.sql.tpl
```

### 模板变量

模板中可以使用以下变量：

- `{{ .Table }}` - 表结构信息
- `{{ .EntityName }}` - 实体名（Pascal）
- `{{ .EntityKebab }}` - 实体名（kebab-case）
- `{{ .Module }}` - 模块名
- `{{ .Operations }}` - 操作列表
- `{{ .Features }}` - 功能开关
- `{{ .HasTree }}` - 是否树形表
- `{{ .HasSoftDelete }}` - 是否软删除

### 模板函数

- `{{ toCamel "user_name" }}` - 转小驼峰
- `{{ toPascal "user_name" }}` - 转大驼峰
- `{{ toKebab "user_name" }}` - 转短横线
- `{{ toSnake "userName" }}` - 转下划线
- `{{ filterListFields .Table.Columns }}` - 过滤列表字段
- `{{ filterFormFields .Table.Columns }}` - 过滤表单字段

## 开发说明

### 项目结构

```
gen/
├── cmd/                    # CLI 命令
├── config/                 # 配置管理
├── engine/                 # 模板渲染引擎
├── generator/              # 生成器核心
├── parser/                 # 数据库解析器
├── types/                  # 类型定义
├── template/               # 模板文件
└── main.go                 # 入口
```

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make fmt
make lint
```

## 依赖

- Go 1.21+
- github.com/gogf/gf/v2
- github.com/go-sql-driver/mysql
- github.com/spf13/cobra
- gopkg.in/yaml.v3

## License

MIT License

## 常见问题

### Q: 如何修改生成字段的表单类型？

A: 修改 `parser/parser.go` 中的 `inferFormType` 函数，添加自定义推断逻辑。

### Q: 如何添加新的模板文件？

A: 在 `template/` 目录下创建新的模板文件，然后在 `generator/generator.go` 中添加生成逻辑。

### Q: 如何支持新的数据库类型？

A: 在 `parser/parser.go` 中添加新的 `parseTableXxx` 方法，并修改 `dataTypeToGo` 和 `dataTypeToTs` 函数。
