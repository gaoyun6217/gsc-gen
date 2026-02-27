package parser

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gfrd/gen/types"
)

// Parser 数据库解析器
type Parser struct {
	db     *sql.DB
	driver string // mysql / postgres
}

// New 创建解析器
func New(dsn string, driver string) (*Parser, error) {
	if driver == "" {
		driver = "mysql"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Parser{
		db:     db,
		driver: driver,
	}, nil
}

// Close 关闭连接
func (p *Parser) Close() error {
	return p.db.Close()
}

// ParseTable 解析表结构
func (p *Parser) ParseTable(ctx context.Context, tableName string) (*types.TableInfo, error) {
	switch p.driver {
	case "mysql":
		return p.parseTableMySQL(ctx, tableName)
	case "postgres":
		return p.parseTablePostgres(ctx, tableName)
	default:
		return p.parseTableMySQL(ctx, tableName)
	}
}

// parseTableMySQL 解析 MySQL 表结构
func (p *Parser) parseTableMySQL(ctx context.Context, tableName string) (*types.TableInfo, error) {
	// 获取表注释
	tableQuery := `
		SELECT TABLE_COMMENT
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`

	var tableComment string
	err := p.db.QueryRowContext(ctx, tableQuery, tableName).Scan(&tableComment)
	if err != nil {
		tableComment = ""
	}

	// 获取列信息
	columnQuery := `
		SELECT
			COLUMN_NAME,
			COLUMN_TYPE,
			DATA_TYPE,
			COLUMN_COMMENT,
			CHARACTER_MAXIMUM_LENGTH,
			NUMERIC_PRECISION,
			NUMERIC_SCALE,
			IS_NULLABLE,
			COLUMN_DEFAULT,
			COLUMN_KEY,
			EXTRA
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`

	rows, err := p.db.QueryContext(ctx, columnQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []*types.ColumnInfo
	var primaryKey string

	for rows.Next() {
		var (
			colName         string
			colType         string
			dataType        string
			comment         string
			maxLength       sql.NullInt64
			precision       sql.NullInt64
			scale           sql.NullInt64
			isNullable      string
			defaultValue    sql.NullString
			columnKey       string
			extra           string
		)

		err := rows.Scan(
			&colName, &colType, &dataType, &comment,
			&maxLength, &precision, &scale,
			&isNullable, &defaultValue, &columnKey, &extra,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}

		// 解析长度和精度
		length := 0
		if maxLength.Valid {
			length = int(maxLength.Int64)
		}
		if length == 0 && precision.Valid {
			length = int(precision.Int64)
		}

		colScale := 0
		if scale.Valid {
			colScale = int(scale.Int64)
		}

		// 解析列类型获取长度 (如 varchar(50))
		if match := regexp.MustCompile(`\((\d+)(?:,(\d+))?\)`).FindStringSubmatch(colType); match != nil {
			fmt.Sscanf(match[1], "%d", &length)
			if match[2] != "" {
				fmt.Sscanf(match[2], "%d", &colScale)
			}
		}

		// 转为 Go 类型
		goType := p.dataTypeToGo(dataType, length, colScale)
		tsType := p.dataTypeToTs(dataType)

		// 判断是否主键
		isPrimary := columnKey == "PRI"
		if isPrimary {
			primaryKey = colName
		}

		// 判断是否自增
		isAutoInc := strings.Contains(extra, "auto_increment")

		// 创建列信息
		col := &types.ColumnInfo{
			Name:         colName,
			NameCamel:    types.ToCamel(colName),
			NamePascal:   types.ToPascal(colName),
			Type:         colType,
			TypeGo:       goType,
			TypeTs:       tsType,
			Comment:      comment,
			Length:       length,
			Precision:    int(precision.Int64),
			Scale:        colScale,
			Nullable:     isNullable == "YES",
			DefaultValue: defaultValue.String,
			IsPrimary:    isPrimary,
			IsAutoInc:    isAutoInc,
			IsListField:  true, // 默认在列表显示
			IsQueryField: false,
			QueryType:    "=",
			FormType:     p.inferFormType(dataType, length, comment),
			Sort:         len(columns),
		}

		// 特殊字段处理
		col.IsListField = p.shouldShowInList(colName, comment)
		col.IsQueryField = p.shouldBeQueryField(colName, comment)

		columns = append(columns, col)
	}

	// 获取索引信息
	indexQuery := `
		SELECT
			INDEX_NAME,
			COLUMN_NAME,
			NON_UNIQUE
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`

	rows, err = p.db.QueryContext(ctx, indexQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	var indexes []*types.IndexInfo
	indexMap := make(map[string]*types.IndexInfo)

	for rows.Next() {
		var (
			indexName  string
			columnName string
			nonUnique  bool
		)

		err := rows.Scan(&indexName, &columnName, &nonUnique)
		if err != nil {
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}

		if idx, ok := indexMap[indexName]; ok {
			idx.Columns = append(idx.Columns, columnName)
		} else {
			indexMap[indexName] = &types.IndexInfo{
				Name:    indexName,
				Columns: []string{columnName},
				Unique:  !nonUnique,
			}
		}
	}

	for _, idx := range indexMap {
		indexes = append(indexes, idx)
	}

	// 判断是否为树形表
	isTreeTable := p.isTreeTable(columns)

	return &types.TableInfo{
		Name:        tableName,
		Comment:     tableComment,
		Columns:     columns,
		PrimaryKey:  primaryKey,
		Indexes:     indexes,
		IsTreeTable: isTreeTable,
	}, nil
}

// parseTablePostgres 解析 PostgreSQL 表结构
func (p *Parser) parseTablePostgres(ctx context.Context, tableName string) (*types.TableInfo, error) {
	// 获取表注释
	tableQuery := `
		SELECT obj_description(c.oid, 'pg_class') AS table_comment
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relname = $1 AND n.nspname = 'public'`

	var tableComment string
	err := p.db.QueryRowContext(ctx, tableQuery, tableName).Scan(&tableComment)
	if err != nil {
		tableComment = ""
	}

	// 获取列信息
	columnQuery := `
		SELECT
			a.attname AS column_name,
			format_type(a.atttypid, a.atttypmod) AS column_type,
			pg_catalog.col_description(a.attrelid, a.attnum) AS comment,
			CASE WHEN a.attnotnull THEN 'NO' ELSE 'YES' END AS is_nullable,
			pg_get_expr(d.adbin, d.adrelid) AS default_value,
			CASE WHEN i.indisprimary THEN 'PRI' ELSE '' END AS column_key
		FROM pg_attribute a
		JOIN pg_class c ON c.oid = a.attrelid
		JOIN pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
		LEFT JOIN pg_index i ON i.indrelid = a.attrelid AND a.attnum = ANY(i.indkey)
		WHERE c.relname = $1 AND n.nspname = 'public' AND a.attnum > 0 AND NOT a.attisdropped
		ORDER BY a.attnum`

	rows, err := p.db.QueryContext(ctx, columnQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []*types.ColumnInfo
	var primaryKey string

	for rows.Next() {
		var (
			colName      string
			colType      string
			comment      string
			isNullable   string
			defaultValue sql.NullString
			columnKey    string
		)

		err := rows.Scan(&colName, &colType, &comment, &isNullable, &defaultValue, &columnKey)
		if err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}

		// 解析基础类型
		dataType := p.extractBaseType(colType)
		length, scale := p.parseTypeModifiers(colType)

		goType := p.dataTypeToGo(dataType, length, scale)
		tsType := p.dataTypeToTs(dataType)

		isPrimary := columnKey == "PRI"
		if isPrimary {
			primaryKey = colName
		}

		col := &types.ColumnInfo{
			Name:         colName,
			NameCamel:    types.ToCamel(colName),
			NamePascal:   types.ToPascal(colName),
			Type:         colType,
			TypeGo:       goType,
			TypeTs:       tsType,
			Comment:      comment,
			Length:       length,
			Nullable:     isNullable == "YES",
			DefaultValue: defaultValue.String,
			IsPrimary:    isPrimary,
			IsListField:  p.shouldShowInList(colName, comment),
			IsQueryField: p.shouldBeQueryField(colName, comment),
			Sort:         len(columns),
		}

		columns = append(columns, col)
	}

	isTreeTable := p.isTreeTable(columns)

	return &types.TableInfo{
		Name:        tableName,
		Comment:     tableComment,
		Columns:     columns,
		PrimaryKey:  primaryKey,
		IsTreeTable: isTreeTable,
	}, nil
}

// dataTypeToGo 数据库类型转 Go 类型
func (p *Parser) dataTypeToGo(dataType string, length int, scale int) string {
	dataType = strings.ToLower(dataType)

	switch {
	case strings.Contains(dataType, "int") && strings.Contains(dataType, "big"):
		return "int64"
	case strings.Contains(dataType, "int") && strings.Contains(dataType, "medium"):
		return "int"
	case strings.Contains(dataType, "int") && strings.Contains(dataType, "small"):
		return "int"
	case strings.Contains(dataType, "int") && strings.Contains(dataType, "tiny"):
		return "int8"
	case strings.Contains(dataType, "int"):
		return "int"
	case strings.Contains(dataType, "bool"):
		return "bool"
	case strings.Contains(dataType, "datetime") || strings.Contains(dataType, "timestamp"):
		return "*gtime.Time"
	case strings.Contains(dataType, "date"):
		return "*gtime.Time"
	case strings.Contains(dataType, "text"):
		return "string"
	case strings.Contains(dataType, "json"):
		return "gjson.RawMessage"
	case strings.Contains(dataType, "decimal") || strings.Contains(dataType, "numeric"):
		return "string"
	case strings.Contains(dataType, "float") || strings.Contains(dataType, "double"):
		return "float64"
	case strings.Contains(dataType, "blob") || strings.Contains(dataType, "binary"):
		return "[]byte"
	default:
		return "string"
	}
}

// dataTypeToTs 数据库类型转 TypeScript 类型
func (p *Parser) dataTypeToTs(dataType string) string {
	dataType = strings.ToLower(dataType)

	switch {
	case strings.Contains(dataType, "int") || strings.Contains(dataType, "decimal") || strings.Contains(dataType, "numeric"):
		return "number"
	case strings.Contains(dataType, "bool"):
		return "boolean"
	case strings.Contains(dataType, "datetime") || strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "date"):
		return "string"
	case strings.Contains(dataType, "json"):
		return "any"
	default:
		return "string"
	}
}

// extractBaseType 提取基础类型 (去除长度等修饰)
func (p *Parser) extractBaseType(colType string) string {
	re := regexp.MustCompile(`^[a-zA-Z]+`)
	match := re.FindString(colType)
	if match != "" {
		return match
	}
	return colType
}

// parseTypeModifiers 解析类型修饰符 (长度、精度)
func (p *Parser) parseTypeModifiers(colType string) (int, int) {
	re := regexp.MustCompile(`\((\d+)(?:,(\d+))?\)`)
	match := re.FindStringSubmatch(colType)
	if len(match) > 1 {
		var length, scale int
		fmt.Sscanf(match[1], "%d", &length)
		if len(match) > 2 && match[2] != "" {
			fmt.Sscanf(match[2], "%d", &scale)
		}
		return length, scale
	}
	return 0, 0
}

// inferFormType 推断表单类型
func (p *Parser) inferFormType(dataType string, length int, comment string) string {
	comment = strings.ToLower(comment)

	// 根据注释推断
	if strings.Contains(comment, "状态") || strings.Contains(comment, "是否") ||
		strings.Contains(comment, "启用") || strings.Contains(comment, "禁用") {
		return "switch"
	}
	if strings.Contains(comment, "性别") || strings.Contains(comment, "类型") {
		return "radio"
	}
	if strings.Contains(comment, "爱好") || strings.Contains(comment, "标签") {
		return "checkbox"
	}
	if strings.Contains(comment, "角色") || strings.Contains(comment, "部门") ||
		strings.Contains(comment, "分类") {
		return "select"
	}
	if strings.Contains(comment, "图片") || strings.Contains(comment, "头像") ||
		strings.Contains(comment, "封面") {
		return "upload"
	}
	if strings.Contains(comment, "时间") || strings.Contains(comment, "日期") {
		if strings.Contains(comment, "开始") || strings.Contains(comment, "结束") {
			return "datetime"
		}
		return "date"
	}
	if strings.Contains(comment, "内容") || strings.Contains(comment, "描述") ||
		strings.Contains(comment, "详情") || strings.Contains(comment, "简介") {
		return "textarea"
	}

	// 根据数据类型推断
	if strings.Contains(dataType, "text") {
		return "textarea"
	}
	if strings.Contains(dataType, "bool") {
		return "switch"
	}
	if strings.Contains(dataType, "date") || strings.Contains(dataType, "time") {
		return "datetime"
	}
	if length > 200 {
		return "textarea"
	}

	return "input"
}

// shouldShowInList 判断是否应该在列表显示
func (p *Parser) shouldShowInList(colName string, comment string) bool {
	// 不显示的字段
	hideFields := []string{"password", "password_hash", "salt", "token", "deleted_at"}
	for _, hide := range hideFields {
		if colName == hide || strings.Contains(strings.ToLower(comment), hide) {
			return false
		}
	}

	// ID 字段通常显示
	if colName == "id" {
		return true
	}

	return true
}

// shouldBeQueryField 判断是否应该作为查询条件
func (p *Parser) shouldBeQueryField(colName string, comment string) bool {
	// 作为查询条件的字段
	queryFields := []string{"name", "username", "code", "status", "type", "email", "phone", "mobile"}
	for _, qf := range queryFields {
		if colName == qf || strings.Contains(strings.ToLower(comment), qf) {
			return true
		}
	}

	return false
}

// isTreeTable 判断是否为树形表
func (p *Parser) isTreeTable(columns []*types.ColumnInfo) bool {
	hasParentID := false
	hasLevel := false
	hasPath := false

	for _, col := range columns {
		name := strings.ToLower(col.Name)
		comment := strings.ToLower(col.Comment)

		if name == "parent_id" || name == "pid" || strings.Contains(comment, "父 id") {
			hasParentID = true
		}
		if name == "level" || strings.Contains(comment, "层级") || strings.Contains(comment, "关系树等级") {
			hasLevel = true
		}
		if name == "path" || strings.Contains(comment, "路径") || strings.Contains(comment, "关系树") {
			hasPath = true
		}
	}

	return hasParentID && (hasLevel || hasPath)
}
