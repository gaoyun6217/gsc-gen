package main

import (
	_ "gf2/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"gf2/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
