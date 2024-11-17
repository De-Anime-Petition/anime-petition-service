package main

import (
	_ "anime_petition/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"anime_petition/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
