package web

import (
	"fmt"
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	//路由注册功能
	AddRoute(method string, path string, handleFunc HandleFunc)
}

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	*router
	mdls []Middleware

	tplEngine TemplateEngine
}

// 缺乏扩展性 如果参数太多，就会该来该去
func NewHTTPServerV1(mdls ...Middleware) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
		mdls:   mdls,
	}
	//errs.

	return res
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
	}

	for _, opt := range opts {
		opt(res)
	}
	fmt.Println("NewHTTPServer", res.tplEngine)
	return res
}

func (h *HTTPServer) Mdls() {
	fmt.Println(h.mdls)
}

func ServerWithTemplateEngine(tplEngine TemplateEngine) HTTPServerOption {

	return func(server *HTTPServer) {
		server.tplEngine = tplEngine

	}
}

func ServerWithMiddleWARE(mdls ...Middleware) HTTPServerOption {

	return func(server *HTTPServer) {
		fmt.Println("ServerWithMiddleWARE invoke")
		server.mdls = mdls
	}
}

func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
	//TODO implement me

	h.router.AddRoute(method, path, handleFunc)
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me

	ctx := &Context{Req: request, Resp: writer, tlpEngine: h.tplEngine}
	root := h.serve
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}
	root = flash2(root)
	root(ctx)
	//flash(ctx)

}

func flash2(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		defer func() {
			ctx.Resp.WriteHeader(ctx.StatusCode)
			ctx.Resp.Write(ctx.RespData)
			return
		}()
		next(ctx)
	}
}

func flash(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.StatusCode)
	ctx.Resp.Write(ctx.RespData)
	return
}

func (h *HTTPServer) serve(ctx *Context) {
	info2, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info2.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	ctx.PathParams = info2.pathParams
	ctx.MatchedRoute = info2.n.path
	/*for _, d := range h.mdls {
		info2.n.handler = d(info2.n.handler)
	}*/
	info2.n.handler(ctx)
}

func (h *HTTPServer) Start(addr string) error {
	/*//TODO implement me
	panic("implement me")*/

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 可以让用户注册after start 回调

	return http.Serve(l, h)
}

func (h *HTTPServer) Start1(addr string) error {
	return http.ListenAndServe(addr, h)
}
