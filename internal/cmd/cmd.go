package cmd

import (
	"anime_petition/internal/controller"
	"anime_petition/middleware"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.SetLogger(g.Log())
			s.BindMiddlewareDefault(middleware.CORSMiddleware, middleware.XSSMiddleware)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(controller.User)
			})
			s.Run()
			return nil
		},
	}
)
