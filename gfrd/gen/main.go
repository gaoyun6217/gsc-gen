package main

import (
	"fmt"
	"os"

	"github.com/gfrd/gen/cmd"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.GetInitCtx()

	// 检查是否有命令行参数，如果没有则进入交互式模式
	if len(os.Args) < 2 {
		// 无参数时进入交互式模式
		if err := cmd.ExecuteInteractive(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 检查是否是交互式命令
	if os.Args[1] == "interactive" || os.Args[1] == "quick" ||
	   os.Args[1] == "history" || os.Args[1] == "rollback" ||
	   os.Args[1] == "import" {
		if err := cmd.ExecuteInteractive(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 否则使用原有的 CLI
	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
