package middleware

import (
	"github.com/gogf/gf/v2/encoding/ghtml"
	"github.com/gogf/gf/v2/net/ghttp"
)

func XSSMiddleware(r *ghttp.Request) {
	r.Request.ParseForm()
	for key, values := range r.Request.Form {
		for i, value := range values {
			r.Request.Form[key][i] = ghtml.SpecialChars(value)
		}
	}
	r.Middleware.Next()
}
