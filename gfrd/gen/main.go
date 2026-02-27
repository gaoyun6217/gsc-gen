package main

import (
	"fmt"
	"os"

	"github.com/gfrd/gen/cmd"
	"github.com/gfrd/gen/web"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.GetInitCtx()

	// 检查是否是 web 命令
	if len(os.Args) > 1 && (os.Args[1] == "web" || os.Args[1] == "serve") {
		if err := web.Serve(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Web server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 无参数或 interactive 命令
	if len(os.Args) < 2 || os.Args[1] == "interactive" || os.Args[1] == "quick" ||
		os.Args[1] == "history" || os.Args[1] == "rollback" || os.Args[1] == "import" {
		if err := cmd.ExecuteInteractive(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 其他 CLI 命令
	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
