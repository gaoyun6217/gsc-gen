package main

import (
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	ctx := gctx.GetInitCtx()

	s := g.Server()

	// 注册路由
	// router.Register(ctx, s)

	s.SetPort(8000)
	s.Run()
}
