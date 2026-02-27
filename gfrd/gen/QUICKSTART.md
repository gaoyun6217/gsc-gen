# GFRD Gen 快速开始指南

## 1. 安装

### 方法 1: 从源码安装

```bash
cd gfrd/gen
go install .
```

### 方法 2: 从 Go module 安装

```bash
go install github.com/gfrd/gen@latest
```

### 方法 3: 直接构建

```bash
cd gfrd/gen
make build
# 可执行文件在 bin/gfrd-gen
```

## 2. 第一个生成示例

### 准备数据库

```sql
CREATE DATABASE IF NOT EXISTS gfrd DEFAULT CHARACTER SET utf8mb4;

USE gfrd;

CREATE TABLE `sys_dept` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '部门 ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父部门 ID',
  `name` varchar(50) NOT NULL COMMENT '部门名称',
  `code` varchar(50) DEFAULT '' COMMENT '部门编码',
  `leader` varchar(20) DEFAULT '' COMMENT '负责人',
  `phone` varchar(11) DEFAULT '' COMMENT '联系电话',
  `email` varchar(100) DEFAULT '' COMMENT '邮箱',
  `status` tinyint DEFAULT '1' COMMENT '状态 (1 正常 0 禁用)',
  `sort` int DEFAULT '0' COMMENT '排序',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='系统部门表';
```

### 生成代码

```bash
# 使用预览模式查看生成的代码
gfrd-gen preview \
  --table="sys_dept" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --module="sys"

# 实际生成代码
gfrd-gen crud \
  --table="sys_dept" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --output="./server" \
  --web-output="./web" \
  --module="sys" \
  --features="add,edit,delete,view,list" \
  --with-test
```

## 3. 生成的代码结构

### 后端文件

```
server/
├── api/sys/
│   └── sys_dept.go              # API 定义
├── internal/handler/sys/
│   └── sys_dept.go              # Handler 实现
├── internal/router/genrouter/
│   └── sys_dept.go              # 路由注册
├── tests/handler/sys/
│   └── sys_dept_test.go         # 单元测试
└── storage/data/generate/
    └── sys_dept_menu.sql        # 菜单 SQL
```

### 前端文件

```
web/src/
├── api/sys/dept/
│   ├── index.ts                 # API 服务
│   └── types.ts                 # TS 类型
└── views/sys/dept/
    ├── index.vue                # 列表页
    └── edit.vue                 # 编辑弹窗
```

## 4. 使用配置文件

创建 `generator.yaml`:

```yaml
database:
  driver: mysql
  dsn: "root:123456@tcp(127.0.0.1:3306)/gfrd?parseTime=true&loc=Local"

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

使用配置文件生成:

```bash
gfrd-gen crud --config=generator.yaml --table="sys_dept"
```

## 5. 命令行参数说明

### 必需参数

| 参数 | 说明 | 示例 |
|------|------|------|
| --table, -t | 表名 | sys_dept |
| --db, -d | 数据库 DSN | mysql:root:123456@tcp(127.0.0.1:3306)/gfrd |

### 可选参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| --output | 后端输出目录 | ./server |
| --web-output | 前端输出目录 | ./web |
| --package | Go 包名 | github.com/gfrd/server |
| --module, -m | 模块名 | sys |
| --features | 要生成的功能 | add,edit,delete,view,list |
| --with-test | 生成单元测试 | false |
| --with-doc | 生成 API 文档 | true |
| --layer-mode | 分层模式 | simple |
| --preview | 仅预览 | false |
| --template | 模板目录 | - |
| --config | 配置文件 | - |

## 6. 功能选项

`--features` 支持以下功能:

| 功能 | 说明 | 生成的代码 |
|------|------|-----------|
| list | 列表查询 | List 方法、搜索表单、数据表格 |
| add | 新增 | Add 方法、新增按钮、新增表单 |
| edit | 修改 | Edit 方法、编辑按钮、编辑表单 |
| delete | 删除 | Delete 方法、删除按钮 |
| view | 详情 | View 方法、详情查看 |
| export | 导出 | 导出按钮、导出逻辑 |
| import | 导入 | 导入按钮、导入逻辑 |

## 7. 智能字段处理

### 自动推断表单类型

```
字段注释包含 → 表单类型
"状态"、"是否" → switch
"性别"、"类型" → radio
"爱好"、"标签" → checkbox
"角色"、"部门" → select
"图片"、"头像" → upload
"时间"、"日期" → date/datetime
"内容"、"描述" → textarea
```

### 自动推断查询字段

```
字段名 → 查询条件
name, username, code → 精确匹配
status, type → 精确匹配
email, phone → 模糊匹配
```

### 自动类型转换

```
数据库类型 → Go 类型 → TS 类型
int/bigint → int64 → number
varchar → string → string
datetime → *gtime.Time → string
tinyint(1) → bool → boolean
```

## 8. 常见问题

### Q: 提示找不到数据库？

A: 检查 DSN 格式是否正确，确保数据库连接信息正确。

```bash
# 测试连接
mysql -h 127.0.0.1 -u root -p123456 -D gfrd
```

### Q: 如何修改生成的字段类型？

A: 修改 `parser/parser.go` 中的 `dataTypeToGo` 方法。

### Q: 如何自定义模板？

A: 复制 `template/` 目录到自定义位置，使用 `--template` 参数指定。

### Q: 生成的代码不符合需求？

A: 可以：
1. 修改模板文件
2. 手动调整生成的代码
3. 提交 issue 或 PR

## 9. 下一步

1. 阅读 [ARCHITECTURE.md](./ARCHITECTURE.md) 了解架构设计
2. 阅读 [README.md](./README.md) 了解完整功能
3. 开始生成你的第一个业务模块！

## 10. 获取帮助

```bash
# 查看帮助
gfrd-gen help

# 查看具体命令帮助
gfrd-gen crud --help
```
