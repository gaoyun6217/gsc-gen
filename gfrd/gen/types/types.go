package types

// TableInfo 表结构信息
type TableInfo struct {
	Name        string        // 表名
	Comment     string        // 表注释
	Columns     []*ColumnInfo // 列信息
	PrimaryKey  string        // 主键列名
	Indexes     []*IndexInfo  // 索引信息
	IsTreeTable bool          // 是否为树形表
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name         string // 列名 (下划线)
	NameCamel    string // 列名 (小驼峰)
	NamePascal   string // 列名 (大驼峰)
	Type         string // 数据库类型
	TypeGo       string // Go 类型
	TypeTs       string // TypeScript 类型
	Comment      string // 注释
	Length       int    // 长度
	Precision    int    // 精度
	Scale        int    // 小数位
	Nullable     bool   // 是否可空
	DefaultValue string // 默认值
	IsPrimary    bool   // 是否主键
	IsAutoInc    bool   // 是否自增
	IsListField  bool   // 是否在列表中显示
	IsQueryField bool   // 是否作为查询条件
	QueryType    string // 查询类型 (=, !=, >, <, LIKE, IN, BETWEEN)
	FormType     string // 表单类型 (input, textarea, select, radio, checkbox, date, datetime, switch, upload)
	DictType     string // 字典类型
	Sort         int    // 排序
}

// IndexInfo 索引信息
type IndexInfo struct {
	Name    string   // 索引名
	Columns []string // 列名
	Unique  bool     // 是否唯一
}

// OperationInfo 操作信息
type OperationInfo struct {
	Name       string // 方法名 (List, Create, Update, Delete, View)
	Comment    string // 注释
	Path       string // API 路径
	Method     string // HTTP 方法
	Tags       string // 标签
	Summary    string // 摘要
	ParamsType string // 参数类型
	RespType   string // 响应类型
	Fields     []*OperationField
}

// OperationField 操作字段
type OperationField struct {
	Name        string // 字段名
	Type        string // 类型
	JsonName    string // JSON 名
	Description string // 描述
	Required    bool   // 是否必填
}

// GeneratorConfig 生成器配置
type GeneratorConfig struct {
	Table      string   // 表名
	DB         string   // 数据库连接
	Output     string   // 后端输出目录
	WebOutput  string   // 前端输出目录
	Package    string   // Go 包名
	Module     string   // 模块名 (sys, org 等)
	Features   []string // 要生成的功能
	WithTest   bool     // 是否生成测试
	WithDoc    bool     // 是否生成文档
	LayerMode  string   // 分层模式 (simple/standard)
	Preview    bool     // 是否仅预览
	Template   string   // 模板目录
	ConfigFile string   // 配置文件路径
}

// RenderData 模板渲染数据
type RenderData struct {
	Table        *TableInfo      // 表信息
	Package      string          // 包名
	Module       string          // 模块名
	EntityName   string          // 实体名 (Pascal)
	EntityKebab  string          // 实体名 (kebab-case)
	EntitySnake  string          // 实体名 (snake_case)
	Operations   []*OperationInfo // 操作列表
	Features     map[string]bool // 功能开关
	HasTree      bool            // 是否有树结构
	HasSoftDelete bool          // 是否有软删除
	HasCreatedAt  bool          // 是否有创建时间
	HasUpdatedAt  bool          // 是否有更新时间
	ImportPackages []string      // 导入的包
}
