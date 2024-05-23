package test

import (
	"dengming20240317/web"
	"dengming20240317/web/session"
	"dengming20240317/web/session/cookie"
	"dengming20240317/web/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	// 非常简单校验
	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store:      memory.NewStore(time.Minute * 15),
		CtxSessKey: "sessKey",
	}
	server := web.NewHTTPServer(web.ServerWithMiddleWARE(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			if ctx.Req.URL.Path == "/login" {
				next(ctx)
				return
			}
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.StatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("重新登录")
				return
			}

			//
			m.RefreshSession(ctx)
			next(ctx)
		}
	},
	))

	server.AddRoute(http.MethodGet, "/user", func(ctx *web.Context) {
		//非login 之前已经登录成功了，所以肯定可以拿出来
		sess, _ := m.GetSession(ctx)
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
		ctx.StatusCode = 200
	})

	server.AddRoute(http.MethodPost, "/login", func(ctx *web.Context) {
		//要先这之前校验用户名和密码

		//校验成功
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespData = []byte("登录失败")
			ctx.StatusCode = http.StatusInternalServerError
			return
		}

		err = sess.Set(ctx.Req.Context(), "nickname", "xiaoming")
		if err != nil {
			ctx.RespData = []byte("登录失败了")
			ctx.StatusCode = http.StatusInternalServerError
			return
		}
		ctx.StatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")
		return
	})

	server.AddRoute(http.MethodPost, "/logout", func(ctx *web.Context) {
		//要先这之前校验用户名和密码

		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.StatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}

		ctx.RespData = []byte("退出登录")
		ctx.StatusCode = http.StatusOK
	})

	server.Start(":8081")
}
