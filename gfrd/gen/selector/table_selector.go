package selector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gfrd/gen/parser"
)

// TableSelector 表选择器
type TableSelector struct {
	parser *parser.Parser
}

// TableOption 表选项
type TableOption struct {
	Name    string
	Comment string
	Selected bool
}

// NewTableSelector 创建表选择器
func NewTableSelector(p *parser.Parser) *TableSelector {
	return &TableSelector{
		parser: p,
	}
}

// ListTables 列出所有表
func (s *TableSelector) ListTables() ([]TableOption, error) {
	tables, err := s.parser.ListTables()
	if err != nil {
		return nil, err
	}

	options := make([]TableOption, 0, len(tables))
	for _, t := range tables {
		options = append(options, TableOption{
			Name:     t.Name,
			Comment:  t.Comment,
			Selected: false,
		})
	}

	return options, nil
}

// InteractiveSelect 交互式选择表
func (s *TableSelector) InteractiveSelect() ([]string, error) {
	options, err := s.ListTables()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		// 显示表列表
		fmt.Println("\n========================================")
		fmt.Println("  选择要生成代码的表 (多选，用逗号分隔)")
		fmt.Println("========================================")
		fmt.Println()

		for i, opt := range options {
			checked := " "
			if opt.Selected {
				checked = "x"
			}
			comment := opt.Comment
			if comment == "" {
				comment = "-"
			}
			fmt.Printf("  [%s] %2d. %-30s (%s)\n", checked, i+1, opt.Name, comment)
		}

		fmt.Println()
		fmt.Println("  输入序号选择/取消，输入 'done' 完成选择，输入 'all' 全选")
		fmt.Print("  > ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "done" || input == "d" {
			break
		}

		if input == "all" {
			for i := range options {
				options[i].Selected = !options[i].Selected
			}
			continue
		}

		// 解析输入的序号
		parts := strings.Split(input, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			idx, err := strconv.Atoi(p)
			if err != nil {
				continue
			}
			if idx >= 1 && idx <= len(options) {
				options[idx-1].Selected = !options[idx-1].Selected
			}
		}
	}

	// 收集选中的表
	selected := make([]string, 0)
	for _, opt := range options {
		if opt.Selected {
			selected = append(selected, opt.Name)
		}
	}

	return selected, nil
}

// SelectSingle 选择单个表
func (s *TableSelector) SelectSingle() (string, error) {
	options, err := s.ListTables()
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n========================================")
	fmt.Println("  选择要生成代码的表")
	fmt.Println("========================================")
	fmt.Println()

	for i, opt := range options {
		comment := opt.Comment
		if comment == "" {
			comment = "-"
		}
		fmt.Printf("  %2d. %-30s (%s)\n", i+1, opt.Name, comment)
	}

	fmt.Println()
	fmt.Printf("  输入序号选择 (1-%d): ", len(options))

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(options) {
		return "", fmt.Errorf("invalid selection")
	}

	return options[idx-1].Name, nil
}
