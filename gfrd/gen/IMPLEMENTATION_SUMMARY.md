# GFRD 交互式 CLI 实现总结

## 实现日期
2026-02-27

## 实现的功能

### 1. 数据库连接与表识别
**文件**: `gen/parser/parser.go`

- 支持 MySQL 和 PostgreSQL 数据库
- 自动解析表结构（字段、类型、注释）
- 智能推断字段类型（表单类型、查询条件）
- 识别树形表结构（parent_id, level, path 字段）
- 新增 `ListTables()` 方法列出所有表

### 2. 交互式表选择器
**文件**: `gen/selector/table_selector.go`

- `ListTables()` - 列出所有表
- `InteractiveSelect()` - 交互式多选表
- `SelectSingle()` - 选择单个表

**功能特点**:
- 显示表名和注释
- 支持多选（输入序号或用逗号分隔）
- 支持全选操作（输入 'all'）
- 完成选择（输入 'done'）

### 3. 生成历史记录
**文件**: `gen/history/history.go`

**数据结构**:
- `GenerationRecord` - 生成记录
  - ID、表名、模块名
  - 生成时间
  - 生成的文件列表
  - 文件内容快照
  - 文件校验和

- `HistoryManager` - 历史记录管理器
  - `AddRecord()` - 添加记录
  - `GetRecords()` - 获取所有记录
  - `GetRecordByID()` - 根据 ID 获取
  - `GetRecordsByTable()` - 根据表名获取
  - `Rollback()` - 回滚到指定记录
  - `DeleteRecord()` - 删除记录
  - `ClearHistory()` - 清空历史

### 4. 交互式 CLI 命令
**文件**: `gen/cmd/interactive_cmd.go`

**命令列表**:

| 命令 | 描述 | 示例 |
|------|------|------|
| `interactive` | 进入交互式模式 | `gfrd-gen interactive` |
| `quick` | 快速生成单个表 | `gfrd-gen quick -t sys_user` |
| `history` | 查看生成历史 | `gfrd-gen history` |
| `rollback` | 回滚到指定版本 | `gfrd-gen rollback -r gen_xxx` |
| `import` | 导入 GoFrame 项目 | `gfrd-gen import -p ./project` |

**交互式流程**:
1. 数据库配置
2. 选择表（多选）
3. 模块配置
4. 输出目录配置
5. 功能选择（list/add/edit/delete/view）
6. 字段配置（列表显示字段、查询条件字段）
7. 逐个表生成代码
8. 保存历史记录

### 5. 生成器增强
**文件**: `gen/generator/generator.go`

**新增导出函数**:
- `PrepareRenderData()` - 准备渲染数据
- `GenerateBackendWithRenderData()` - 生成后端代码
- `GenerateFrontendWithRenderData()` - 生成前端代码

**辅助函数**:
- `removePrefixForConfig()` - 移除表前缀
- `buildOperationsForConfig()` - 构建操作列表
- `buildFeaturesForConfig()` - 构建功能开关
- `hasSoftDeleteInTable()` - 检查软删除字段
- `hasCreatedAtInTable()` - 检查创建时间字段
- `hasUpdatedAtInTable()` - 检查更新时间字段

### 6. 入口点更新
**文件**: `gen/main.go`

- 无参数时自动进入交互式模式
- 自动识别交互式命令并路由
- 保持原有 CLI 命令兼容

## 文件结构

```
gfrd/gen/
├── cmd/
│   ├── cmd.go                  # 原有 CLI 命令
│   └── interactive_cmd.go      # 新增：交互式 CLI 命令
├── config/
│   └── config.go               # 配置管理
├── engine/
│   └── renderer.go             # 模板渲染引擎
├── generator/
│   └── generator.go            # 代码生成器（已增强）
├── history/
│   └── history.go              # 新增：历史记录管理
├── parser/
│   └── parser.go               # 数据库解析器（已增强）
├── selector/
│   └── table_selector.go       # 新增：表选择器
├── types/
│   ├── types.go                # 类型定义
│   └── namecase.go             # 名称转换
├── template/
│   ├── backend/                # 后端模板
│   ├── frontend/               # 前端模板
│   └── sql/                    # SQL 模板
├── main.go                     # 入口（已更新）
└── INTERACTIVE_GUIDE.md        # 新增：使用指南
```

## 使用示例

### 进入交互式模式

```bash
cd gfrd/gen
go run main.go
```

### 快速生成

```bash
go run main.go quick \
  --table="sys_user" \
  --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
  --module="sys"
```

### 查看历史

```bash
go run main.go history
```

### 回滚

```bash
go run main.go rollback --record-id=gen_1708992000000000000
```

### 导入项目

```bash
go run main.go import \
  --project-path="/path/to/goframe/project"
```

## 智能推断规则

### 表单类型推断

| 条件 | 推断类型 |
|------|----------|
| 注释包含"状态"/"是否" | switch |
| 注释包含"性别"/"类型" | radio |
| 注释包含"爱好"/"标签" | checkbox |
| 注释包含"角色"/"部门" | select |
| 注释包含"图片"/"头像" | upload |
| 注释包含"时间"/"日期" | datetime |
| 注释包含"内容"/"描述" | textarea |
| 数据类型为 text | textarea |
| 数据类型为 bool | switch |
| 长度 > 200 | textarea |

### 查询字段推断

以下字段名自动添加为查询条件：
- `name`, `username`, `code`, `status`, `type`, `email`, `phone`, `mobile`

### 列表字段推断

以下字段不在列表显示：
- `password`, `password_hash`, `salt`, `token`, `deleted_at`

## 生成记录管理

### 存储位置
`./.gen_history/history.json`

### 记录内容
- 表信息（表名、注释、字段数）
- 生成配置（模块、目录、功能）
- 文件列表（路径、类型、内容、校验和）
- 时间戳

### 回滚机制
1. 查看历史记录
2. 选择要回滚的记录 ID
3. 确认操作（需要输入 y/yes）
4. 恢复所有文件到快照状态

## 兼容性

### 数据库
- MySQL 5.7+
- PostgreSQL 14+

### Go 版本
- Go 1.18+

### GoFrame 版本
- GoFrame 2.7.1+

## 后续可扩展功能

1. **Web UI 界面** - 基于浏览器的交互式生成
2. **模板市场** - 分享和下载社区模板
3. **批量导出** - 导出所有表结构配置
4. **差异对比** - 对比表结构变化生成迁移代码
5. **API 文档生成** - 自动生成 Swagger/OpenAPI 文档
6. **单元测试生成** - 生成 Handler 和 Service 测试

## 技术文档

- [INTERACTIVE_GUIDE.md](INTERACTIVE_GUIDE.md) - 详细使用指南
- [README.md](README.md) - 项目说明
- [QUICKSTART.md](QUICKSTART.md) - 快速开始

## 测试建议

1. 测试 MySQL 数据库连接
2. 测试 PostgreSQL 数据库连接
3. 测试交互式表选择
4. 测试字段配置功能
5. 测试生成历史记录
6. 测试回滚操作
7. 测试 GoFrame 项目导入

## 已知限制

1. 当前仅支持单数据库操作
2. 表关联关系需要手动配置
3. 复杂查询条件需要手动调整
4. 树形表格的前端组件需要额外实现

## 提交建议

```bash
cd gfrd/gen
git add .
git commit -m "feat: 实现交互式 CLI 与生成历史记录功能

- 新增交互式表选择器
- 支持字段级别配置
- 实现生成历史记录管理
- 支持版本回滚
- 支持 GoFrame 项目导入
- 新增 interactive/quick/history/rollback/import 命令
"
```

---

**实现者**: AI Assistant
**完成时间**: 2026-02-27
