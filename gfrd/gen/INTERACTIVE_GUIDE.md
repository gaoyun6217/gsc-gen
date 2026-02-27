# GFRD 交互式代码生成器使用指南

## 概述

GFRD 交互式代码生成器是一个基于 GoFrame 2 的全栈代码生成工具，提供交互式 CLI 界面，支持数据库表选择、字段配置、生成历史记录和回滚功能。

## 快速开始

### 1. 启动交互式模式

```bash
cd gfrd/gen
go run main.go interactive
```

或者直接运行（无参数时自动进入交互式模式）：

```bash
go run main.go
```

### 2. 导入现有的 GoFrame 项目

```bash
go run main.go import --project-path=/path/to/goframe/project
```

### 3. 快速生成单个表

```bash
go run main.go quick \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --module="sys"
```

### 4. 查看生成历史

```bash
go run main.go history
```

### 5. 回滚到指定版本

```bash
go run main.go rollback --record-id=gen_1234567890
```

## 功能特性

### 交互式表选择

- 列出数据库中所有表
- 支持多选（输入序号或用逗号分隔）
- 显示表注释
- 支持全选操作

### 字段级别配置

- 选择哪些字段在列表页面显示
- 选择哪些字段作为查询条件
- 支持批量切换

### 功能选择

可以选择生成以下功能：
1. 列表查询 (list)
2. 新增 (add)
3. 修改 (edit)
4. 删除 (delete)
5. 详情查看 (view)

### 生成历史记录

每次生成都会保存记录，包括：
- 表名和模块名
- 生成时间
- 生成的文件列表
- 文件内容快照（用于回滚）
- 文件校验和

### 回滚支持

可以随时回滚到之前的任意生成版本：
- 查看历史记录
- 选择要回滚的版本
- 确认后恢复文件

## 命令行参数

### interactive 命令

```bash
gfrd-gen interactive [flags]

Flags:
  -f, --configure-fields   是否配置字段 (默认 true)
  -h, --help               帮助信息
```

### quick 命令

```bash
gfrd-gen quick [flags]

Flags:
  -t, --table string       表名 (必需)
  -d, --db string          数据库连接 (必需)
  -m, --module string      模块名 (默认 "sys")
      --output string      后端输出目录 (默认 "./server")
      --web string         前端输出目录 (默认 "./web")
```

### history 命令

```bash
gfrd-gen history
```

### rollback 命令

```bash
gfrd-gen rollback [flags]

Flags:
  -r, --record-id string   记录 ID (必需)
```

### import 命令

```bash
gfrd-gen import [flags]

Flags:
  -p, --project-path string   GoFrame 项目路径 (必需)
```

## 交互式流程示例

### 步骤 1: 数据库配置

```
--- 数据库配置 ---

  数据库连接 (mysql:root:123456@tcp(127.0.0.1:3306)/gfrd):
  数据库连接成功!
```

### 步骤 2: 选择表

```
========================================
  选择要生成代码的表 (多选，用逗号分隔)
========================================

  [ ]  1. sys_dept                      (系统部门表)
  [ ]  2. sys_user                      (系统用户表)
  [ ]  3. sys_role                      (系统角色表)

  输入序号选择/取消，输入 'done' 完成选择，输入 'all' 全选
  >
```

### 步骤 3: 模块配置

```
  模块名 (sys):
```

### 步骤 4: 输出目录

```
  后端输出目录 (./server):
  前端输出目录 (./web):
```

### 步骤 5: 功能选择

```
--- 选择要生成的功能 ---

  [1] 列表查询 (list)
  [2] 新增 (add)
  [3] 修改 (edit)
  [4] 删除 (delete)
  [5] 详情查看 (view)

  输入序号选择 (默认全选，用逗号分隔):
```

### 步骤 6: 字段配置

```
--- 配置表 sys_dept 的字段 ---

列表显示字段:
  [x]  1. name                (部门名称)
  [x]  2. code                (部门编码)
  [x]  3. leader              (负责人)
  ...

  输入序号切换显示状态 (逗号分隔，直接回车跳过):
```

## 生成的文件

### 后端文件

- `api/{module}/{entity}.go` - API 定义
- `internal/handler/{module}/{entity}.go` - Handler 实现
- `internal/router/genrouter/{entity}.go` - 路由注册
- `storage/data/generate/{entity}_menu.sql` - 菜单 SQL

### 前端文件

- `src/api/{module}/{entity}/index.ts` - API 服务
- `src/api/{module}/{entity}/types.ts` - TypeScript 类型
- `src/views/{module}/{entity}/index.vue` - 列表页面
- `src/views/{module}/{entity}/edit.vue` - 编辑弹窗

## 数据库支持

支持以下数据库：

- MySQL 5.7+
- PostgreSQL 14+

### MySQL 连接示例

```
mysql:username:password@tcp(host:port)/database
```

### PostgreSQL 连接示例

```
postgres:username:password@tcp(host:port)/database
```

## 智能字段推断

### 表单类型推断

| 字段名/类型 | 推断的表单类型 |
|-------------|----------------|
| `status` | switch |
| `remark`, `description` | textarea |
| `email` | input (邮箱) |
| `phone`, `mobile` | input (电话) |
| `password` | input (密码) |
| `date`, `time` | datetime |
| `tinyint(1)` | switch |

### 查询字段推断

以下字段名会自动添加为查询条件：
- `name` - 名称
- `code` - 编码
- `status` - 状态
- `type` - 类型
- `email` - 邮箱
- `phone` - 电话

## 历史记录管理

### 存储位置

历史记录存储在 `./.gen_history/history.json`

### 记录格式

```json
{
  "id": "gen_1708992000000000000",
  "table": "sys_dept",
  "module": "sys",
  "generated_at": "2026-02-27T10:00:00Z",
  "files": [...],
  "table_comment": "系统部门表",
  "field_count": 12,
  "config": {...},
  "checksum": "abc123..."
}
```

## 常见问题

### Q: 如何更改默认的表前缀？

修改 `generator.go` 中的 `removePrefixForConfig` 函数，添加或修改前缀列表。

### Q: 如何自定义模板？

将自定义模板放入 `template` 目录，与内置模板保持相同的目录结构。

### Q: 回滚操作安全吗？

回滚前会提示确认，建议先提交到 git 再进行回滚操作。

### Q: 支持批量生成吗？

支持。在交互式模式中选择多个表后，会依次处理每个表。

## 技术栈

- **后端**: GoFrame 2.7.1+, Go 1.18+
- **前端**: Vue 3, TypeScript, NaiveUI
- **数据库**: MySQL, PostgreSQL
- **CLI**: Cobra

## 更新日志

### v1.0.0 (2026-02-27)

- 初始版本
- 支持交互式表选择
- 支持字段级别配置
- 支持生成历史记录
- 支持版本回滚
- 支持 GoFrame 项目导入
