package web

import (
	"context"
	"strings"
	"time"

	"github.com/gfrd/gen/generator"
	"github.com/gfrd/gen/history"
	"github.com/gfrd/gen/parser"
	"github.com/gfrd/gen/types"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

// Handler Web 处理器
type Handler struct{}

// DBReq 数据库请求
type DBReq struct {
	DSN  string `json:"dsn" v:"required"`
	Type string `json:"type" d:"mysql"`
}

// TestDB 测试数据库连接
func (h *Handler) TestDB(r *ghttp.Request) {
	var req DBReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	p, err := parser.New(req.DSN, req.Type)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "连接失败：" + err.Error(),
		})
		return
	}
	defer p.Close()

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "连接成功",
	})
}

// ListTables 获取表列表
func (h *Handler) ListTables(r *ghttp.Request) {
	var req DBReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	p, err := parser.New(req.DSN, req.Type)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer p.Close()

	tables, err := p.ListTables()
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	result := make([]g.Map, 0, len(tables))
	for _, t := range tables {
		result = append(result, g.Map{
			"name":    t.Name,
			"comment": t.Comment,
		})
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    result,
	})
}

// GetTableDetail 获取表详情
func (h *Handler) GetTableDetail(r *ghttp.Request) {
	var req struct {
		DBReq
		Table string `json:"table" v:"required"`
	}
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	p, err := parser.New(req.DSN, req.Type)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer p.Close()

	table, err := p.ParseTable(context.Background(), req.Table)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 转换为前端格式
	columns := make([]g.Map, 0, len(table.Columns))
	for _, col := range table.Columns {
		columns = append(columns, g.Map{
			"name":         col.Name,
			"nameCamel":    col.NameCamel,
			"namePascal":   col.NamePascal,
			"type":         col.Type,
			"typeGo":       col.TypeGo,
			"comment":      col.Comment,
			"length":       col.Length,
			"nullable":     col.Nullable,
			"isPrimary":    col.IsPrimary,
			"isAutoInc":    col.IsAutoInc,
			"isListField":  col.IsListField,
			"isQueryField": col.IsQueryField,
			"queryType":    col.QueryType,
			"formType":     col.FormType,
		})
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"data": g.Map{
			"name":       table.Name,
			"comment":    table.Comment,
			"columns":    columns,
			"primaryKey": table.PrimaryKey,
			"isTree":     table.IsTreeTable,
		},
	})
}

// GenerateReq 生成请求
type GenerateReq struct {
	DSN      string   `json:"dsn" v:"required"`
	Type     string   `json:"type" d:"mysql"`
	Tables   []string `json:"tables" v:"required"`
	Module   string   `json:"module" d:"sys"`
	Output   string   `json:"output" d:"./server"`
	Web      string   `json:"web" d:"./web"`
	Features []string `json:"features"`
}

// Generate 生成代码
func (h *Handler) Generate(r *ghttp.Request) {
	var req GenerateReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if len(req.Features) == 0 {
		req.Features = []string{"list", "add", "edit", "delete", "view"}
	}

	ctx := r.Context()

	// 逐个表生成
	for _, table := range req.Tables {
		cfg := &generator.Config{
			Table:     table,
			DB:        req.DSN,
			Output:    req.Output,
			WebOutput: req.Web,
			Module:    req.Module,
			Features:  req.Features,
		}

		gen := generator.NewGenerator(cfg)
		if err := gen.Generate(ctx); err != nil {
			r.Response.WriteJson(g.Map{
				"success": false,
				"message": "生成失败：" + err.Error(),
			})
			return
		}
	}

	// 保存到历史记录
	historyManager, _ := history.NewHistoryManager("./.gen_history")
	var lastRecordID string
	if historyManager != nil {
		for _, table := range req.Tables {
			record := &history.GenerationRecord{
				ID:          history.GenerateRecordID(),
				Table:       table,
				Module:      req.Module,
				GeneratedAt: time.Now(),
				Config: history.GeneratorConfig{
					Output:    req.Output,
					WebOutput: req.Web,
					Features:  strings.Join(req.Features, ","),
				},
			}
			historyManager.AddRecord(record)
			lastRecordID = record.ID
		}
	}

	r.Response.WriteJson(g.Map{
		"success":  true,
		"message":  "生成成功",
		"recordId": lastRecordID,
	})
}

// PreviewReq 预览请求
type PreviewReq struct {
	DSN      string   `json:"dsn" v:"required"`
	Type     string   `json:"type" d:"mysql"`
	Table    string   `json:"table" v:"required"`
	Module   string   `json:"module" d:"sys"`
	Features []string `json:"features"`
}

// Preview 预览代码
func (h *Handler) Preview(r *ghttp.Request) {
	var req PreviewReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if len(req.Features) == 0 {
		req.Features = []string{"list", "add", "edit", "delete", "view"}
	}

	ctx := r.Context()

	p, err := parser.New(req.DSN, req.Type)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer p.Close()

	table, err := p.ParseTable(ctx, req.Table)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 生成预览内容
	files := make([]g.Map, 0)

	// 后端 API 预览
	apiFile := "api/" + req.Module + "/" + gstr.ToLower(types.ToPascal(table.Name)) + ".go"
	files = append(files, g.Map{
		"path":    apiFile,
		"type":    "backend",
		"content": "预览功能开发中",
	})

	r.Response.WriteJson(g.Map{
		"success": true,
		"files":   files,
	})
}

// Download 下载代码
func (h *Handler) Download(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"success": false,
		"message": "功能开发中",
	})
}

// ListHistory 获取历史记录
func (h *Handler) ListHistory(r *ghttp.Request) {
	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	records := historyManager.GetRecords()
	result := make([]g.Map, 0, len(records))
	for _, rec := range records {
		result = append(result, g.Map{
			"id":           rec.ID,
			"table":        rec.Table,
			"module":       rec.Module,
			"generatedAt":  rec.GeneratedAt.Format("2006-01-02 15:04:05"),
			"tableComment": rec.TableComment,
			"fieldCount":   rec.FieldCount,
			"fileCount":    len(rec.Files),
		})
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"data":    result,
	})
}

// GetHistoryDetail 获取历史详情
func (h *Handler) GetHistoryDetail(r *ghttp.Request) {
	recordID := r.Get("id").String()

	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	record := historyManager.GetRecordByID(recordID)
	if record == nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "记录不存在",
		})
		return
	}

	files := make([]g.Map, 0, len(record.Files))
	for _, f := range record.Files {
		files = append(files, g.Map{
			"path":    f.Path,
			"type":    f.Type,
			"content": f.Content,
		})
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"data": g.Map{
			"id":           record.ID,
			"table":        record.Table,
			"module":       record.Module,
			"generatedAt":  record.GeneratedAt.Format("2006-01-02 15:04:05"),
			"tableComment": record.TableComment,
			"fieldCount":   record.FieldCount,
			"files":        files,
		},
	})
}

// RollbackReq 回滚请求
type RollbackReq struct {
	RecordID string `json:"recordId" v:"required"`
}

// Rollback 回滚
func (h *Handler) Rollback(r *ghttp.Request) {
	var req RollbackReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := historyManager.Rollback(req.RecordID); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "回滚成功",
	})
}

// DeleteHistory 删除历史记录
func (h *Handler) DeleteHistory(r *ghttp.Request) {
	recordID := r.Get("id").String()

	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := historyManager.DeleteRecord(recordID); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "删除成功",
	})
}

// toKebab 转换为短横线命名
func toKebab(s string) string {
	result := gstr.ReplaceByMap(s, map[string]string{
		"A": "-a", "B": "-b", "C": "-c", "D": "-d", "E": "-e",
		"F": "-f", "G": "-g", "H": "-h", "I": "-i", "J": "-j",
		"K": "-k", "L": "-l", "M": "-m", "N": "-n", "O": "-o",
		"P": "-p", "Q": "-q", "R": "-r", "S": "-s", "T": "-t",
		"U": "-u", "V": "-v", "W": "-w", "X": "-x", "Y": "-y", "Z": "-z",
	})
	return strings.TrimLeft(result, "-")
}
