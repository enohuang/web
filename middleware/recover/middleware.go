package recover

import "dengming20240317/web"

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespData = m.Data
					ctx.StatusCode = m.StatusCode
				}
			}()

			next(ctx)
		}

	}
}
