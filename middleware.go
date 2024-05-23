package web

// Middleware 函数是的责任链模式， 洋葱模式
type Middleware func(next HandleFunc) HandleFunc

type MiddlewareV1 interface {
	Invoke(next HandleFunc) HandleFunc
}

// 拦截器模式
type Interceptor interface {
	Before(ctx *Context)
	After(ctx *Context)
	Interfacer(ctx *Context)
}
