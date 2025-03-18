package middleware

import (
	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 允许跨域请求中间件
func CORSMiddleware(r *ghttp.Request) {
	ctx := gctx.New()
	corsOptions := r.Response.DefaultCORSOptions()
	value, _ := g.Cfg().Get(ctx, "server.AllowDomain")
	corsOptions.AllowDomain = value.Strings()
	corsOptions.AllowMethods = "GET, POST"
	// g.Log().Infof(ctx, "AllowDomain: %#v", corsOptions.AllowDomain)
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}
