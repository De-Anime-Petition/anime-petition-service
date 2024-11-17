package cmd

import (
	"anime_petition/internal/controller"
	"context"

	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

// 允许跨域请求中间件
func Middleware(r *ghttp.Request) {
	ctx := gctx.New()
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = g.Cfg().MustGet(ctx, "server.AllowDomain").Strings()
	// g.Log().Infof(ctx, "AllowDomain: %#v", corsOptions.AllowDomain)
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.BindMiddlewareDefault(Middleware)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(controller.User)
			})
			s.Run()
			return nil
		},
	}
)
