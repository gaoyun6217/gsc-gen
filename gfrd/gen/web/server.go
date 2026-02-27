package web

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Serve 启动 Web 服务器
func Serve(ctx context.Context) error {
	s := g.Server("gfrd-web")

	// 配置
	s.SetPort(8199)
	s.SetServerRoot("./web/static")
	s.SetDumpRouterMap(false)

	// 中间件
	s.Use(func(r *ghttp.Request) {
		r.Response.Header().Set("Access-Control-Allow-Origin", "*")
		r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			r.Response.Status = 200
			return
		}
		r.Middleware.Next()
	})

	// API 路由
	s.Group("/api", func(group *ghttp.RouterGroup) {
		// 数据库相关 API
		group.POST("/db/test", new(Handler).TestDB)
		group.POST("/db/tables", new(Handler).ListTables)
		group.POST("/db/table/detail", new(Handler).GetTableDetail)

		// 代码生成 API
		group.POST("/generate", new(Handler).Generate)
		group.POST("/generate/preview", new(Handler).Preview)
		group.GET("/generate/download", new(Handler).Download)

		// 历史记录 API
		group.GET("/history", new(Handler).ListHistory)
		group.GET("/history/:id", new(Handler).GetHistoryDetail)
		group.POST("/history/rollback", new(Handler).Rollback)
		group.DELETE("/history/:id", new(Handler).DeleteHistory)
	})

	// 前端页面 - 所有路由返回 index.html
	s.BindHandler("/", func(r *ghttp.Request) {
		content, err := ioutil.ReadFile("./web/static/index.html")
		if err != nil {
			r.Response.WriteStatus(404)
			return
		}
		r.Response.Write(content)
	})

	fmt.Println("========================================")
	fmt.Println("  GFRD Web Admin")
	fmt.Println("  http://localhost:8199")
	fmt.Println("========================================")

	s.Run()
	return nil
}
