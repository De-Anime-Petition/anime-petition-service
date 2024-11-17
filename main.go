package main

import (
	_ "anime_petition/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"anime_petition/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
