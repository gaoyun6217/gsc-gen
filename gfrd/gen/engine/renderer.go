package engine

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gfrd/gen/types"
)

// Renderer 模板渲染器
type Renderer struct {
	templateDir string
	funcMap     template.FuncMap
}

// NewRenderer 创建渲染器
func NewRenderer(templateDir string) *Renderer {
	if templateDir == "" {
		templateDir = "template"
	}

	r := &Renderer{
		templateDir: templateDir,
		funcMap:     template.FuncMap{},
	}

	// 注册模板函数
	r.registerFuncs()

	return r
}

// registerFuncs 注册模板函数
func (r *Renderer) registerFuncs() {
	r.funcMap = template.FuncMap{
		// 名称转换
		"toCamel":  types.ToCamel,
		"toPascal": types.ToPascal,
		"toKebab":  types.ToKebab,
		"toSnake":  types.ToSnake,
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"title":    strings.Title,
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"split":    strings.Split,
		"join":     strings.Join,
		"lenInt":   func(i int) int { return i },
		"printf":   fmt.Sprintf,
		"print":    fmt.Sprint,
		"println":  fmt.Sprintln,

		// 自定义函数
		"importTime":     importTime,
		"importGFrame":   importGFrame,
		"buildTags":      buildTags,
		"buildPermissions": buildPermissions,
		"filterQueryFields": filterQueryFields,
		"filterListFields":  filterListFields,
		"filterFormFields":  filterFormFields,
	}
}

// Render 渲染模板
func (r *Renderer) Render(ctx context.Context, tmplFile string, data *types.RenderData) (string, error) {
	tmplPath := filepath.Join(r.templateDir, tmplFile)

	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New(filepath.Base(tmplFile)).Funcs(r.funcMap).Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderAndWrite 渲染并写入文件
func (r *Renderer) RenderAndWrite(ctx context.Context, tmplFile string, outputPath string, data *types.RenderData) error {
	content, err := r.Render(ctx, tmplFile, data)
	if err != nil {
		return err
	}

	// 创建目录
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// 模板辅助函数

func importTime(data *types.RenderData) bool {
	for _, col := range data.Table.Columns {
		if col.TypeGo == "*gtime.Time" {
			return true
		}
	}
	return false
}

func importGFrame(data *types.RenderData) bool {
	return true // 总是导入 gframe
}

func buildTags(data *types.RenderData) string {
	tags := []string{"系统管理"}
	if data.Module != "" {
		tags = []string{data.Module}
	}
	return strings.Join(tags, ",")
}

func buildPermissions(data *types.RenderData, op string) string {
	base := fmt.Sprintf("/%s/%s", data.Module, data.EntityKebab)
	switch op {
	case "list":
		return fmt.Sprintf(`["%s/list"]`, base)
	case "add":
		return fmt.Sprintf(`["%s/add"]`, base)
	case "edit":
		return fmt.Sprintf(`["%s/edit"]`, base)
	case "delete":
		return fmt.Sprintf(`["%s/delete"]`, base)
	case "view":
		return fmt.Sprintf(`["%s/view"]`, base)
	}
	return "[]"
}

func filterQueryFields(columns []*types.ColumnInfo) []*types.ColumnInfo {
	var result []*types.ColumnInfo
	for _, col := range columns {
		if col.IsQueryField {
			result = append(result, col)
		}
	}
	return result
}

func filterListFields(columns []*types.ColumnInfo) []*types.ColumnInfo {
	var result []*types.ColumnInfo
	for _, col := range columns {
		if col.IsListField && !col.IsPrimary {
			result = append(result, col)
		}
	}
	// 确保 id 在最前面
	if len(result) > 0 && result[0].Name != "id" {
		for i, col := range result {
			if col.Name == "id" {
				// 将 id 移到最前面
				result = append([]*types.ColumnInfo{col}, append(result[:i], result[i+1:]...)...)
				break
			}
		}
	}
	return result
}

func filterFormFields(columns []*types.ColumnInfo) []*types.ColumnInfo {
	var result []*types.ColumnInfo
	hideFields := []string{"id", "created_at", "updated_at", "deleted_at", "password_hash", "salt"}

	for _, col := range columns {
		skip := false
		for _, hide := range hideFields {
			if col.Name == hide {
				skip = true
				break
			}
		}
		if !skip {
			result = append(result, col)
		}
	}
	return result
}
