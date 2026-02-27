package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gfrd/gen/config"
	"github.com/gfrd/gen/engine"
	"github.com/gfrd/gen/generator"
	"github.com/gfrd/gen/history"
	"github.com/gfrd/gen/parser"
	"github.com/gfrd/gen/selector"
	"github.com/gfrd/gen/types"
	"github.com/spf13/cobra"
)

// InteractiveConfig 交互式配置
type InteractiveConfig struct {
	DB            string
	Module        string
	Output        string
	WebOutput     string
	Tables        []string
	Features      []string
	ConfigureFields bool
}

// ExecuteInteractive 执行交互式 CLI
func ExecuteInteractive(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Use:   "gfrd-gen",
		Short: "GFRD Interactive Code Generator",
		Long: `GFRD 交互式代码生成器 - 基于 GoFrame 2 的全栈代码生成工具

功能特性:
  - 交互式表选择
  - 字段级别配置
  - 生成历史记录
  - 支持回滚操作

使用示例:
  gfrd-gen interactive  # 进入交互式模式
  gfrd-gen quick        # 快速生成模式
  gfrd-gen history      # 查看生成历史
  gfrd-gen rollback     # 回滚到指定版本
`,
	}

	// 添加子命令
	rootCmd.AddCommand(cmdInteractive())
	rootCmd.AddCommand(cmdQuick())
	rootCmd.AddCommand(cmdHistory())
	rootCmd.AddCommand(cmdRollback())
	rootCmd.AddCommand(cmdImportProject())

	return rootCmd.ExecuteContext(ctx)
}

// cmdInteractive 交互式生成命令
func cmdInteractive() *cobra.Command {
	var configureFields bool

	cmd := &cobra.Command{
		Use:   "interactive",
		Short: "进入交互式代码生成模式",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return runInteractive(ctx, configureFields)
		},
	}

	cmd.Flags().BoolVarP(&configureFields, "configure-fields", "f", true, "是否配置字段")
	return cmd
}

// cmdQuick 快速生成命令
func cmdQuick() *cobra.Command {
	var (
		table  string
		db     string
		module string
		output string
		web    string
	)

	cmd := &cobra.Command{
		Use:   "quick",
		Short: "快速生成单个表的代码",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &generator.Config{
				Table:     table,
				DB:        db,
				Output:    output,
				WebOutput: web,
				Module:    module,
			}
			return generator.NewGenerator(cfg).Generate(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&table, "table", "t", "", "表名")
	cmd.Flags().StringVarP(&db, "db", "d", "", "数据库连接")
	cmd.Flags().StringVarP(&module, "module", "m", "sys", "模块名")
	cmd.Flags().StringVar(&output, "output", "./server", "后端输出目录")
	cmd.Flags().StringVar(&web, "web", "./web", "前端输出目录")

	return cmd
}

// cmdHistory 查看历史命令
func cmdHistory() *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "查看代码生成历史",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showHistory()
		},
	}
}

// cmdRollback 回滚命令
func cmdRollback() *cobra.Command {
	var recordID string

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "回滚到指定的生成版本",
		RunE: func(cmd *cobra.Command, args []string) error {
			if recordID == "" {
				return fmt.Errorf("--record-id is required")
			}
			return doRollback(recordID)
		},
	}

	cmd.Flags().StringVarP(&recordID, "record-id", "r", "", "记录 ID")

	return cmd
}

// cmdImportProject 导入 GoFrame 项目
func cmdImportProject() *cobra.Command {
	var projectPath string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "导入现有的 GoFrame 项目",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectPath == "" {
				return fmt.Errorf("--project-path is required")
			}
			return importProject(projectPath)
		},
	}

	cmd.Flags().StringVarP(&projectPath, "project-path", "p", "", "GoFrame 项目路径")

	return cmd
}

// runInteractive 运行交互式模式
func runInteractive(ctx context.Context, configureFields bool) error {
	fmt.Println("\n========================================")
	fmt.Println("  GFRD 交互式代码生成器")
	fmt.Println("========================================")

	// 1. 数据库配置
	db, err := configureDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// 2. 选择表
	tableSelector := selector.NewTableSelector(db)
	tables, err := tableSelector.InteractiveSelect()
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		fmt.Println("未选择任何表")
		return nil
	}

	// 3. 模块配置
	module := configureModule()

	// 4. 配置输出目录
	outputDir, webOutputDir := configureOutput()

	// 5. 选择功能
	features := selectFeatures()

	// 6. 逐个表配置字段并生成
	for _, table := range tables {
		fmt.Printf("\n正在处理表：%s\n", table)

		// 解析表结构
		tableInfo, err := db.ParseTable(ctx, table)
		if err != nil {
			fmt.Printf("解析表失败：%v\n", err)
			continue
		}

		// 配置字段
		if configureFields {
			tableInfo = configureTableFields(tableInfo)
		}

		// 生成代码
		cfg := &generator.Config{
			Table:     table,
			DB:        "", // 已解析
			Output:    outputDir,
			WebOutput: webOutputDir,
			Module:    module,
			Features:  features,
		}

		// 使用已解析的表信息生成
		if err := generateWithTableInfo(ctx, cfg, tableInfo); err != nil {
			fmt.Printf("生成失败：%v\n", err)
			continue
		}

		fmt.Printf("表 %s 处理完成!\n", table)
	}

	fmt.Println("\n========================================")
	fmt.Println("  代码生成完成!")
	fmt.Println("========================================")

	return nil
}

// configureDatabase 配置数据库
func configureDatabase() (*parser.Parser, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n--- 数据库配置 ---")
	fmt.Println()

	// 数据库连接
	fmt.Print("  数据库连接 (mysql:root:123456@tcp(127.0.0.1:3306)/gfrd): ")
	dsn, _ := reader.ReadString('\n')
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		dsn = "mysql:root:123456@tcp(127.0.0.1:3306)/gfrd"
	}

	// 解析 DSN
	dbType := "mysql"
	if strings.HasPrefix(dsn, "mysql:") {
		dbType = "mysql"
		dsn = strings.TrimPrefix(dsn, "mysql:")
	} else if strings.HasPrefix(dsn, "postgres:") {
		dbType = "postgres"
		dsn = strings.TrimPrefix(dsn, "postgres:")
	}

	p, err := parser.New(dsn, dbType)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败：%w", err)
	}

	fmt.Println("数据库连接成功!")
	return p, nil
}

// configureModule 配置模块名
func configureModule() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Print("  模块名 (sys): ")
	module, _ := reader.ReadString('\n')
	module = strings.TrimSpace(module)
	if module == "" {
		module = "sys"
	}

	return module
}

// configureOutput 配置输出目录
func configureOutput() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Print("  后端输出目录 (./server): ")
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)
	if output == "" {
		output = "./server"
	}

	fmt.Print("  前端输出目录 (./web): ")
	webOutput, _ := reader.ReadString('\n')
	webOutput = strings.TrimSpace(webOutput)
	if webOutput == "" {
		webOutput = "./web"
	}

	return output, webOutput
}

// selectFeatures 选择功能
func selectFeatures() []string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("--- 选择要生成的功能 ---")
	fmt.Println()
	fmt.Println("  [1] 列表查询 (list)")
	fmt.Println("  [2] 新增 (add)")
	fmt.Println("  [3] 修改 (edit)")
	fmt.Println("  [4] 删除 (delete)")
	fmt.Println("  [5] 详情查看 (view)")
	fmt.Println()
	fmt.Print("  输入序号选择 (默认全选，用逗号分隔): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return []string{"list", "add", "edit", "delete", "view"}
	}

	features := make([]string, 0)
	featureMap := map[string]string{
		"1": "list",
		"2": "add",
		"3": "edit",
		"4": "delete",
		"5": "view",
	}

	parts := strings.Split(input, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if f, ok := featureMap[p]; ok {
			features = append(features, f)
		}
	}

	return features
}

// configureTableFields 配置表字段
func configureTableFields(table *types.TableInfo) *types.TableInfo {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Printf("--- 配置表 %s 的字段 ---\n", table.Name)
	fmt.Println()

	// 配置列表显示字段
	fmt.Println("列表显示字段:")
	for i, col := range table.Columns {
		if col.Name == "id" || col.Name == "created_at" || col.Name == "updated_at" {
			continue
		}
		checked := " "
		if col.IsListField {
			checked = "x"
		}
		fmt.Printf("  [%s] %2d. %-20s (%s)\n", checked, i+1, col.Name, col.Comment)
	}

	fmt.Println()
	fmt.Print("  输入序号切换显示状态 (逗号分隔，直接回车跳过): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		parts := strings.Split(input, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			idx, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			if idx >= 1 && idx <= len(table.Columns) {
				table.Columns[idx-1].IsListField = !table.Columns[idx-1].IsListField
			}
		}
	}

	// 配置查询字段
	fmt.Println()
	fmt.Println("查询条件字段:")
	for i, col := range table.Columns {
		if !col.IsQueryField {
			continue
		}
		fmt.Printf("  [x] %2d. %-20s (%s)\n", i+1, col.Name, col.Comment)
	}

	fmt.Println()
	fmt.Print("  输入序号添加/删除查询条件 (逗号分隔，直接回车跳过): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		parts := strings.Split(input, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			idx, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			if idx >= 1 && idx <= len(table.Columns) {
				col := table.Columns[idx-1]
				col.IsQueryField = !col.IsQueryField
				if col.IsQueryField {
					col.QueryType = "="
				}
			}
		}
	}

	return table
}

// generateWithTableInfo 使用已解析的表信息生成代码
func generateWithTableInfo(ctx context.Context, cfg *generator.Config, table *types.TableInfo) error {
	// 准备渲染数据
	renderData := generator.PrepareRenderData(cfg, table)

	// 渲染模板
	renderer := engine.NewRenderer("")

	// 创建历史记录
	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		fmt.Printf("创建历史管理器失败：%v\n", err)
	}

	record := &history.GenerationRecord{
		ID:           history.GenerateRecordID(),
		Table:        table.Name,
		Module:       cfg.Module,
		GeneratedAt:  time.Now(),
		TableComment: table.Comment,
		FieldCount:   len(table.Columns),
		Config: history.GeneratorConfig{
			Output:    cfg.Output,
			WebOutput: cfg.WebOutput,
			Package:   cfg.Package,
			Features:  strings.Join(cfg.Features, ","),
		},
		Files: make([]history.GeneratedFile, 0),
	}

	// 生成后端代码
	backendFiles, err := generator.GenerateBackendWithRenderData(renderer, cfg.Output, renderData)
	if err != nil {
		return err
	}

	// 生成前端代码
	frontendFiles, err := generator.GenerateFrontendWithRenderData(renderer, cfg.WebOutput, renderData)
	if err != nil {
		return err
	}

	// 记录生成的文件
	for _, f := range backendFiles {
		content, _ := os.ReadFile(f.Path)
		record.Files = append(record.Files, history.GeneratedFile{
			Path:      f.Path,
			Type:      "backend",
			Content:   string(content),
			Checksum:  history.CalculateChecksum(string(content)),
			CreatedAt: time.Now(),
		})
	}

	for _, f := range frontendFiles {
		content, _ := os.ReadFile(f.Path)
		record.Files = append(record.Files, history.GeneratedFile{
			Path:      f.Path,
			Type:      "frontend",
			Content:   string(content),
			Checksum:  history.CalculateChecksum(string(content)),
			CreatedAt: time.Now(),
		})
	}

	// 保存历史记录
	if historyManager != nil {
		if err := historyManager.AddRecord(record); err != nil {
			fmt.Printf("保存历史记录失败：%v\n", err)
		}
	}

	return nil
}

// showHistory 显示历史记录
func showHistory() error {
	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		return err
	}

	records := historyManager.GetRecords()

	if len(records) == 0 {
		fmt.Println("暂无生成历史")
		return nil
	}

	fmt.Println("\n========================================")
	fmt.Println("  代码生成历史")
	fmt.Println("========================================")
	fmt.Println()

	for i, r := range records {
		fmt.Printf("%3d. [%s] %s\n", i+1, r.ID, r.Table)
		fmt.Printf("    模块：%s | 字段数：%d | 时间：%s\n", r.Module, r.FieldCount, r.GeneratedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("    文件数：%d\n", len(r.Files))
		fmt.Println()
	}

	return nil
}

// doRollback 执行回滚
func doRollback(recordID string) error {
	historyManager, err := history.NewHistoryManager("./.gen_history")
	if err != nil {
		return err
	}

	record := historyManager.GetRecordByID(recordID)
	if record == nil {
		return fmt.Errorf("记录不存在：%s", recordID)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("确定要回滚到 %s 的版本吗？此操作将覆盖当前文件。[y/N]: ", record.Table)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("已取消回滚")
		return nil
	}

	if err := historyManager.Rollback(recordID); err != nil {
		return fmt.Errorf("回滚失败：%w", err)
	}

	fmt.Printf("回滚成功！已恢复到 %s 的版本。\n", record.Table)
	return nil
}

// importProject 导入 GoFrame 项目
func importProject(projectPath string) error {
	fmt.Println("\n========================================")
	fmt.Println("  导入 GoFrame 项目")
	fmt.Println("========================================")
	fmt.Println()

	// 检查项目结构
	requiredDirs := []string{"api", "internal", "manifest"}
	missingDirs := make([]string, 0)

	for _, dir := range requiredDirs {
		path := filepath.Join(projectPath, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missingDirs = append(missingDirs, dir)
		}
	}

	if len(missingDirs) > 0 {
		return fmt.Errorf("不是有效的 GoFrame 项目，缺少目录：%v", missingDirs)
	}

	// 读取配置文件
	configPath := filepath.Join(projectPath, "manifest", "config", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Printf("读取配置文件失败：%v\n", err)
		} else {
			fmt.Printf("项目配置加载成功!\n")
			fmt.Printf("  数据库：%s\n", cfg.Database.Driver)
			fmt.Printf("  DSN: %s\n", cfg.Database.DSN)
		}
	}

	// 保存项目路径配置
	configData, _ := json.MarshalIndent(map[string]string{
		"project_path": projectPath,
		"imported_at":  time.Now().Format(time.RFC3339),
	}, "", "  ")

	projectConfigFile := filepath.Join(projectPath, ".gen_project.json")
	os.WriteFile(projectConfigFile, configData, 0644)

	fmt.Println()
	fmt.Println("项目导入成功!")
	fmt.Println()
	fmt.Println("接下来可以使用以下命令:")
	fmt.Println("  gfrd-gen interactive  # 进入交互式生成模式")
	fmt.Println("  gfrd-gen quick -t <表名>  # 快速生成单个表")

	return nil
}
