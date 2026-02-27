package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gfrd/gen/engine"
	"github.com/gfrd/gen/parser"
	"github.com/gfrd/gen/types"
)

// GeneratedFile 已生成的文件信息
type GeneratedFile struct {
	Path string
	Type string // backend/frontend/sql
}

// Config 生成器配置
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
	Template     string   // 模板目录
	OnlyBackend  bool     // 仅生成后端
	OnlyFrontend bool     // 仅生成前端
}

// Generator 代码生成器
type Generator struct {
	cfg      *Config
	parser   *parser.Parser
	renderer *engine.Renderer
}

// NewGenerator 创建生成器
func NewGenerator(cfg *Config) *Generator {
	return &Generator{
		cfg:      cfg,
		renderer: engine.NewRenderer(cfg.Template),
	}
}

// Generate 执行生成
func (g *Generator) Generate(ctx context.Context) error {
	// 初始化数据库解析器
	dbType, dsn := g.parseDBConfig(g.cfg.DB)
	p, err := parser.New(dsn, dbType)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer p.Close()

	g.parser = p

	// 解析表结构
	table, err := p.ParseTable(ctx, g.cfg.Table)
	if err != nil {
		return fmt.Errorf("failed to parse table: %w", err)
	}

	// 准备渲染数据
	renderData := g.prepareRenderData(table)

	// 生成代码
	if !g.cfg.OnlyFrontend {
		if err := g.generateBackend(ctx, renderData); err != nil {
			return err
		}
	}

	if !g.cfg.OnlyBackend {
		if err := g.generateFrontend(ctx, renderData); err != nil {
			return err
		}
	}

	return nil
}

// parseDBConfig 解析数据库配置
func (g *Generator) parseDBConfig(dsn string) (string, string) {
	if strings.HasPrefix(dsn, "mysql:") {
		return "mysql", strings.TrimPrefix(dsn, "mysql:")
	}
	if strings.HasPrefix(dsn, "postgres:") {
		return "postgres", strings.TrimPrefix(dsn, "postgres:")
	}
	// 默认 MySQL
	return "mysql", dsn
}

// prepareRenderData 准备渲染数据
func (g *Generator) prepareRenderData(table *types.TableInfo) *types.RenderData {
	entityName := types.ToPascal(g.removePrefix(table.Name))

	return &types.RenderData{
		Table:        table,
		Package:      g.cfg.Module,
		Module:       g.cfg.Module,
		EntityName:   entityName,
		EntityKebab:  strings.ToLower(types.ToKebab(entityName)),
		EntitySnake:  types.ToSnake(entityName),
		Operations:   g.buildOperations(table),
		Features:     g.buildFeatures(),
		HasTree:      table.IsTreeTable,
		HasSoftDelete: g.hasSoftDelete(table),
		HasCreatedAt: g.hasCreatedAt(table),
		HasUpdatedAt: g.hasUpdatedAt(table),
	}
}

// removePrefix 移除表前缀
func (g *Generator) removePrefix(tableName string) string {
	prefixes := []string{"sys_", "admin_", "hg_", "t_", "tb_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(tableName, prefix) {
			return strings.TrimPrefix(tableName, prefix)
		}
	}
	return tableName
}

// buildOperations 构建操作列表
func (g *Generator) buildOperations(table *types.TableInfo) []*types.OperationInfo {
	var ops []*types.OperationInfo
	features := g.buildFeatures()

	if features["list"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "List",
			Comment: "获取" + table.Comment + "列表",
			Path:    "/" + g.cfg.Module + "/" + types.ToKebab(types.ToPascal(g.removePrefix(table.Name))) + "/list",
			Method:  "get",
			Tags:    g.cfg.Module,
			Summary: "获取" + table.Comment + "列表",
		})
	}

	if features["add"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "Add",
			Comment: "添加" + table.Comment,
			Path:    "/" + g.cfg.Module + "/" + types.ToKebab(types.ToPascal(g.removePrefix(table.Name))) + "/add",
			Method:  "post",
			Tags:    g.cfg.Module,
			Summary: "添加" + table.Comment,
		})
	}

	if features["edit"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "Edit",
			Comment: "修改" + table.Comment,
			Path:    "/" + g.cfg.Module + "/" + types.ToKebab(types.ToPascal(g.removePrefix(table.Name))) + "/edit",
			Method:  "post",
			Tags:    g.cfg.Module,
			Summary: "修改" + table.Comment,
		})
	}

	if features["delete"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "Delete",
			Comment: "删除" + table.Comment,
			Path:    "/" + g.cfg.Module + "/" + types.ToKebab(types.ToPascal(g.removePrefix(table.Name))) + "/delete",
			Method:  "post",
			Tags:    g.cfg.Module,
			Summary: "删除" + table.Comment,
		})
	}

	if features["view"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "View",
			Comment: "查看" + table.Comment + "详情",
			Path:    "/" + g.cfg.Module + "/" + types.ToKebab(types.ToPascal(g.removePrefix(table.Name))) + "/view",
			Method:  "get",
			Tags:    g.cfg.Module,
			Summary: "查看" + table.Comment + "详情",
		})
	}

	return ops
}

// buildFeatures 构建功能开关
func (g *Generator) buildFeatures() map[string]bool {
	features := make(map[string]bool)
	for _, f := range g.cfg.Features {
		features[strings.TrimSpace(f)] = true
	}
	// 默认至少包含 list
	if len(features) == 0 {
		features["list"] = true
	}
	return features
}

// hasSoftDelete 是否有软删除字段
func (g *Generator) hasSoftDelete(table *types.TableInfo) bool {
	for _, col := range table.Columns {
		if col.Name == "deleted_at" {
			return true
		}
	}
	return false
}

// hasCreatedAt 是否有创建时间字段
func (g *Generator) hasCreatedAt(table *types.TableInfo) bool {
	for _, col := range table.Columns {
		if col.Name == "created_at" {
			return true
		}
	}
	return false
}

// hasUpdatedAt 是否有更新时间字段
func (g *Generator) hasUpdatedAt(table *types.TableInfo) bool {
	for _, col := range table.Columns {
		if col.Name == "updated_at" {
			return true
		}
	}
	return false
}

// generateBackend 生成后端代码
func (g *Generator) generateBackend(ctx context.Context, data *types.RenderData) error {
	fmt.Printf("Generating backend code for %s...\n", data.EntityName)

	basePath := g.cfg.Output
	modulePath := filepath.Join(basePath, "internal", "handler", g.cfg.Module)
	apiPath := filepath.Join(basePath, "api", g.cfg.Module)

	// 创建目录
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(apiPath, 0755); err != nil {
		return err
	}

	// 生成 API 文件
	apiFile := filepath.Join(apiPath, data.EntitySnake+".go")
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "backend/api.go.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + apiFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "backend/api.go.tpl", apiFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", apiFile)
	}

	// 生成 Handler 文件
	handlerFile := filepath.Join(modulePath, data.EntitySnake+".go")
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "backend/handler.go.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + handlerFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "backend/handler.go.tpl", handlerFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", handlerFile)
	}

	// 生成 Service 文件 (standard 模式)
	if g.cfg.LayerMode == "standard" {
		serviceFile := filepath.Join(basePath, "internal", "service", g.cfg.Module, data.EntitySnake+".go")
		if err := os.MkdirAll(filepath.Dir(serviceFile), 0755); err != nil {
			return err
		}
		if g.cfg.Preview {
			content, err := g.renderer.Render(ctx, "backend/service.go.tpl", data)
			if err != nil {
				return err
			}
			fmt.Println("=== " + serviceFile + " ===")
			fmt.Println(content)
			fmt.Println()
		} else {
			if err := g.renderer.RenderAndWrite(ctx, "backend/service.go.tpl", serviceFile, data); err != nil {
				return err
			}
			fmt.Printf("  Created: %s\n", serviceFile)
		}
	}

	// 生成路由文件
	routerFile := filepath.Join(basePath, "internal", "router", "genrouter", data.EntitySnake+".go")
	if err := os.MkdirAll(filepath.Dir(routerFile), 0755); err != nil {
		return err
	}
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "backend/router.go.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + routerFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "backend/router.go.tpl", routerFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", routerFile)
	}

	// 生成 SQL 菜单文件
	sqlFile := filepath.Join(basePath, "storage", "data", "generate", data.EntitySnake+"_menu.sql")
	if err := os.MkdirAll(filepath.Dir(sqlFile), 0755); err != nil {
		return err
	}
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "sql/menu.sql.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + sqlFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "sql/menu.sql.tpl", sqlFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", sqlFile)
	}

	// 生成测试文件
	if g.cfg.WithTest {
		testFile := filepath.Join(basePath, "tests", "handler", g.cfg.Module, data.EntitySnake+"_test.go")
		if err := os.MkdirAll(filepath.Dir(testFile), 0755); err != nil {
			return err
		}
		if g.cfg.Preview {
			content, err := g.renderer.Render(ctx, "backend/test.go.tpl", data)
			if err != nil {
				return err
			}
			fmt.Println("=== " + testFile + " ===")
			fmt.Println(content)
			fmt.Println()
		} else {
			if err := g.renderer.RenderAndWrite(ctx, "backend/test.go.tpl", testFile, data); err != nil {
				return err
			}
			fmt.Printf("  Created: %s\n", testFile)
		}
	}

	return nil
}

// generateFrontend 生成前端代码
func (g *Generator) generateFrontend(ctx context.Context, data *types.RenderData) error {
	fmt.Printf("Generating frontend code for %s...\n", data.EntityName)

	basePath := g.cfg.WebOutput
	apiPath := filepath.Join(basePath, "api", g.cfg.Module, data.EntityKebab)
	viewPath := filepath.Join(basePath, "views", g.cfg.Module, data.EntityKebab)

	// 创建目录
	if err := os.MkdirAll(apiPath, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(viewPath, 0755); err != nil {
		return err
	}

	// 生成 API 文件
	apiFile := filepath.Join(apiPath, "index.ts")
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "frontend/api.ts.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + apiFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "frontend/api.ts.tpl", apiFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", apiFile)
	}

	// 生成 TypeScript 类型文件
	typesFile := filepath.Join(apiPath, "types.ts")
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "frontend/types.ts.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + typesFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "frontend/types.ts.tpl", typesFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", typesFile)
	}

	// 生成列表页
	indexFile := filepath.Join(viewPath, "index.vue")
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "frontend/index.vue.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + indexFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "frontend/index.vue.tpl", indexFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", indexFile)
	}

	// 生成编辑弹窗
	editFile := filepath.Join(viewPath, "edit.vue")
	if g.cfg.Preview {
		content, err := g.renderer.Render(ctx, "frontend/edit.vue.tpl", data)
		if err != nil {
			return err
		}
		fmt.Println("=== " + editFile + " ===")
		fmt.Println(content)
		fmt.Println()
	} else {
		if err := g.renderer.RenderAndWrite(ctx, "frontend/edit.vue.tpl", editFile, data); err != nil {
			return err
		}
		fmt.Printf("  Created: %s\n", editFile)
	}

	return nil
}

// PrepareRenderData 准备渲染数据（导出函数）
func PrepareRenderData(cfg *Config, table *types.TableInfo) *types.RenderData {
	entityName := types.ToPascal(removePrefixForConfig(cfg, table.Name))

	return &types.RenderData{
		Table:        table,
		Package:      cfg.Module,
		Module:       cfg.Module,
		EntityName:   entityName,
		EntityKebab:  strings.ToLower(types.ToKebab(entityName)),
		EntitySnake:  types.ToSnake(entityName),
		Operations:   buildOperationsForConfig(cfg, table),
		Features:     buildFeaturesForConfig(cfg.Features),
		HasTree:      table.IsTreeTable,
		HasSoftDelete: hasSoftDeleteInTable(table),
		HasCreatedAt: hasCreatedAtInTable(table),
		HasUpdatedAt: hasUpdatedAtInTable(table),
	}
}

// GenerateBackendWithRenderData 使用渲染数据生成后端代码（导出函数）
func GenerateBackendWithRenderData(renderer *engine.Renderer, output string, data *types.RenderData) ([]GeneratedFile, error) {
	files := make([]GeneratedFile, 0)
	ctx := context.Background()

	basePath := output
	modulePath := filepath.Join(basePath, "internal", "handler", data.Module)
	apiPath := filepath.Join(basePath, "api", data.Module)
	routerPath := filepath.Join(basePath, "internal", "router", "genrouter")
	sqlPath := filepath.Join(basePath, "storage", "data", "generate")

	// 创建目录
	dirs := []string{modulePath, apiPath, routerPath, sqlPath}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return files, err
		}
	}

	// 生成 API 文件
	apiFile := filepath.Join(apiPath, data.EntitySnake+".go")
	if err := renderer.RenderAndWrite(ctx, "backend/api.go.tpl", apiFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: apiFile, Type: "backend"})

	// 生成 Handler 文件
	handlerFile := filepath.Join(modulePath, data.EntitySnake+".go")
	if err := renderer.RenderAndWrite(ctx, "backend/handler.go.tpl", handlerFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: handlerFile, Type: "backend"})

	// 生成路由文件
	routerFile := filepath.Join(routerPath, data.EntitySnake+".go")
	if err := renderer.RenderAndWrite(ctx, "backend/router.go.tpl", routerFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: routerFile, Type: "backend"})

	// 生成 SQL 菜单文件
	sqlFile := filepath.Join(sqlPath, data.EntitySnake+"_menu.sql")
	if err := renderer.RenderAndWrite(ctx, "sql/menu.sql.tpl", sqlFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: sqlFile, Type: "sql"})

	return files, nil
}

// GenerateFrontendWithRenderData 使用渲染数据生成前端代码（导出函数）
func GenerateFrontendWithRenderData(renderer *engine.Renderer, webOutput string, data *types.RenderData) ([]GeneratedFile, error) {
	files := make([]GeneratedFile, 0)
	ctx := context.Background()

	basePath := webOutput
	apiPath := filepath.Join(basePath, "api", data.Module, data.EntityKebab)
	viewPath := filepath.Join(basePath, "views", data.Module, data.EntityKebab)

	// 创建目录
	if err := os.MkdirAll(apiPath, 0755); err != nil {
		return files, err
	}
	if err := os.MkdirAll(viewPath, 0755); err != nil {
		return files, err
	}

	// 生成 API 文件
	apiFile := filepath.Join(apiPath, "index.ts")
	if err := renderer.RenderAndWrite(ctx, "frontend/api.ts.tpl", apiFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: apiFile, Type: "frontend"})

	// 生成 TypeScript 类型文件
	typesFile := filepath.Join(apiPath, "types.ts")
	if err := renderer.RenderAndWrite(ctx, "frontend/types.ts.tpl", typesFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: typesFile, Type: "frontend"})

	// 生成列表页
	indexFile := filepath.Join(viewPath, "index.vue")
	if err := renderer.RenderAndWrite(ctx, "frontend/index.vue.tpl", indexFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: indexFile, Type: "frontend"})

	// 生成编辑弹窗
	editFile := filepath.Join(viewPath, "edit.vue")
	if err := renderer.RenderAndWrite(ctx, "frontend/edit.vue.tpl", editFile, data); err != nil {
		return files, err
	}
	files = append(files, GeneratedFile{Path: editFile, Type: "frontend"})

	return files, nil
}

// 辅助函数（包内使用）
func removePrefixForConfig(cfg *Config, tableName string) string {
	prefixes := []string{"sys_", "admin_", "hg_", "t_", "tb_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(tableName, prefix) {
			return strings.TrimPrefix(tableName, prefix)
		}
	}
	return tableName
}

func buildOperationsForConfig(cfg *Config, table *types.TableInfo) []*types.OperationInfo {
	var ops []*types.OperationInfo
	features := buildFeaturesForConfig(cfg.Features)

	module := cfg.Module
	entityName := removePrefixForConfig(cfg, table.Name)
	entityKebab := types.ToKebab(types.ToPascal(entityName))

	if features["list"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "List",
			Comment: "获取" + table.Comment + "列表",
			Path:    "/" + module + "/" + entityKebab + "/list",
			Method:  "get",
			Tags:    module,
			Summary: "获取" + table.Comment + "列表",
		})
	}

	if features["add"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "Add",
			Comment: "添加" + table.Comment,
			Path:    "/" + module + "/" + entityKebab + "/add",
			Method:  "post",
			Tags:    module,
			Summary: "添加" + table.Comment,
		})
	}

	if features["edit"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "Edit",
			Comment: "修改" + table.Comment,
			Path:    "/" + module + "/" + entityKebab + "/edit",
			Method:  "post",
			Tags:    module,
			Summary: "修改" + table.Comment,
		})
	}

	if features["delete"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "Delete",
			Comment: "删除" + table.Comment,
			Path:    "/" + module + "/" + entityKebab + "/delete",
			Method:  "post",
			Tags:    module,
			Summary: "删除" + table.Comment,
		})
	}

	if features["view"] {
		ops = append(ops, &types.OperationInfo{
			Name:    "View",
			Comment: "查看" + table.Comment + "详情",
			Path:    "/" + module + "/" + entityKebab + "/view",
			Method:  "get",
			Tags:    module,
			Summary: "查看" + table.Comment + "详情",
		})
	}

	return ops
}

func buildFeaturesForConfig(features []string) map[string]bool {
	featureMap := make(map[string]bool)
	for _, f := range features {
		featureMap[strings.TrimSpace(f)] = true
	}
	if len(featureMap) == 0 {
		featureMap["list"] = true
	}
	return featureMap
}

func hasSoftDeleteInTable(table *types.TableInfo) bool {
	for _, col := range table.Columns {
		if col.Name == "deleted_at" {
			return true
		}
	}
	return false
}

func hasCreatedAtInTable(table *types.TableInfo) bool {
	for _, col := range table.Columns {
		if col.Name == "created_at" {
			return true
		}
	}
	return false
}

func hasUpdatedAtInTable(table *types.TableInfo) bool {
	for _, col := range table.Columns {
		if col.Name == "updated_at" {
			return true
		}
	}
	return false
}
