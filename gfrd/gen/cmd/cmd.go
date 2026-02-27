package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gfrd/gen/generator"
	"github.com/spf13/cobra"
)

// Execute 执行 CLI
func Execute(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Use:   "gfrd-gen",
		Short: "GFRD Code Generator - 全栈代码生成器",
		Long: `GFRD Code Generator 是基于 GoFrame 2 和 SoybeanAdmin 的全栈代码生成器。

支持生成:
  - 后端：API 定义、Handler、Service、DAO、Entity
  - 前端：API 服务、TypeScript 类型、Vue 组件
  - SQL:  菜单权限脚本

使用示例:
  gfrd-gen crud --table="sys_user" --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd"
  gfrd-gen backend --table="sys_user" --output="./server"
  gfrd-gen frontend --table="sys_user" --web-output="./web"
  gfrd-gen preview --table="sys_user" --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd"
`,
	}

	// 添加子命令
	rootCmd.AddCommand(genCrudCmd())
	rootCmd.AddCommand(genBackendCmd())
	rootCmd.AddCommand(genFrontendCmd())
	rootCmd.AddCommand(genPreviewCmd())

	return rootCmd.ExecuteContext(ctx)
}

// genCrudCmd 生成完整 CRUD 命令
func genCrudCmd() *cobra.Command {
	var cfg Config

	cmd := &cobra.Command{
		Use:   "crud",
		Short: "Generate full-stack CRUD code",
		Long: `生成完整的 CRUD 代码，包括后端 API、Handler、Service 和前端 Vue 组件

示例:
  gfrd-gen crud \
    --table="sys_user" \
    --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
    --output="./server" \
    --web-output="./web" \
    --module="sys" \
    --features="add,edit,delete,view,list" \
    --with-test \
    --with-doc
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if cfg.Table == "" {
				return fmt.Errorf("--table is required")
			}
			if cfg.DB == "" {
				return fmt.Errorf("--db is required")
			}

			gen := generator.NewGenerator(&generator.Config{
				Table:      cfg.Table,
				DB:         cfg.DB,
				Output:     cfg.Output,
				WebOutput:  cfg.WebOutput,
				Package:    cfg.Package,
				Module:     cfg.Module,
				Features:   strings.Split(cfg.Features, ","),
				WithTest:   cfg.WithTest,
				WithDoc:    cfg.WithDoc,
				LayerMode:  cfg.LayerMode,
				Preview:    cfg.Preview,
				Template:   cfg.Template,
			})

			return gen.Generate(ctx)
		},
	}

	cmd.Flags().StringVarP(&cfg.Table, "table", "t", "", "Table name")
	cmd.Flags().StringVarP(&cfg.DB, "db", "d", "", "Database DSN")
	cmd.Flags().StringVar(&cfg.Output, "output", "./server", "Backend output directory")
	cmd.Flags().StringVar(&cfg.WebOutput, "web-output", "./web", "Frontend output directory")
	cmd.Flags().StringVar(&cfg.Package, "package", "github.com/gfrd/server", "Go package name")
	cmd.Flags().StringVarP(&cfg.Module, "module", "m", "sys", "Module name")
	cmd.Flags().StringVar(&cfg.Features, "features", "add,edit,delete,view,list", "Features to generate")
	cmd.Flags().BoolVar(&cfg.WithTest, "with-test", false, "Generate unit tests")
	cmd.Flags().BoolVar(&cfg.WithDoc, "with-doc", true, "Generate API documentation")
	cmd.Flags().StringVar(&cfg.LayerMode, "layer-mode", "simple", "Layer mode: simple/standard")
	cmd.Flags().BoolVar(&cfg.Preview, "preview", false, "Preview generated code")
	cmd.Flags().StringVar(&cfg.Template, "template", "", "Template directory")

	return cmd
}

// genBackendCmd 仅生成后端命令
func genBackendCmd() *cobra.Command {
	var cfg Config

	cmd := &cobra.Command{
		Use:   "backend",
		Short: "Generate backend code only",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if cfg.Table == "" {
				return fmt.Errorf("--table is required")
			}
			if cfg.DB == "" {
				return fmt.Errorf("--db is required")
			}

			gen := generator.NewGenerator(&generator.Config{
				Table:       cfg.Table,
				DB:          cfg.DB,
				Output:      cfg.Output,
				Package:     cfg.Package,
				Module:      cfg.Module,
				WithTest:    cfg.WithTest,
				WithDoc:     cfg.WithDoc,
				LayerMode:   "simple",
				OnlyBackend: true,
			})

			return gen.Generate(ctx)
		},
	}

	cmd.Flags().StringVarP(&cfg.Table, "table", "t", "", "Table name")
	cmd.Flags().StringVarP(&cfg.DB, "db", "d", "", "Database DSN")
	cmd.Flags().StringVar(&cfg.Output, "output", "./server", "Backend output directory")
	cmd.Flags().StringVar(&cfg.Package, "package", "github.com/gfrd/server", "Go package name")
	cmd.Flags().StringVarP(&cfg.Module, "module", "m", "sys", "Module name")
	cmd.Flags().BoolVar(&cfg.WithTest, "with-test", false, "Generate unit tests")
	cmd.Flags().BoolVar(&cfg.WithDoc, "with-doc", true, "Generate API documentation")

	return cmd
}

// genFrontendCmd 仅生成前端命令
func genFrontendCmd() *cobra.Command {
	var cfg Config

	cmd := &cobra.Command{
		Use:   "frontend",
		Short: "Generate frontend code only",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if cfg.Table == "" {
				return fmt.Errorf("--table is required")
			}
			if cfg.DB == "" {
				return fmt.Errorf("--db is required")
			}

			gen := generator.NewGenerator(&generator.Config{
				Table:        cfg.Table,
				DB:           cfg.DB,
				WebOutput:    cfg.WebOutput,
				Module:       cfg.Module,
				OnlyFrontend: true,
			})

			return gen.Generate(ctx)
		},
	}

	cmd.Flags().StringVarP(&cfg.Table, "table", "t", "", "Table name")
	cmd.Flags().StringVarP(&cfg.DB, "db", "d", "", "Database DSN")
	cmd.Flags().StringVar(&cfg.WebOutput, "web-output", "./web", "Frontend output directory")
	cmd.Flags().StringVarP(&cfg.Module, "module", "m", "sys", "Module name")

	return cmd
}

// genPreviewCmd 预览生成结果命令
func genPreviewCmd() *cobra.Command {
	var cfg Config

	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview generated code without writing files",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if cfg.Table == "" {
				return fmt.Errorf("--table is required")
			}
			if cfg.DB == "" {
				return fmt.Errorf("--db is required")
			}

			gen := generator.NewGenerator(&generator.Config{
				Table:   cfg.Table,
				DB:      cfg.DB,
				Module:  cfg.Module,
				Preview: true,
			})

			return gen.Generate(ctx)
		},
	}

	cmd.Flags().StringVarP(&cfg.Table, "table", "t", "", "Table name")
	cmd.Flags().StringVarP(&cfg.DB, "db", "d", "", "Database DSN")
	cmd.Flags().StringVarP(&cfg.Module, "module", "m", "sys", "Module name")

	return cmd
}

// Config CLI 配置结构
type Config struct {
	Table      string
	DB         string
	Output     string
	WebOutput  string
	Package    string
	Module     string
	Features   string
	WithTest   bool
	WithDoc    bool
	LayerMode  string
	Preview    bool
	Template   string
	OnlyBackend  bool
	OnlyFrontend bool
}
