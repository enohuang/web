package errhdl

import (
	"dengming20240317/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`<html><body><h1>Not Found</h1></body></html>`))
	server := web.NewHTTPServer(web.ServerWithMiddleWARE(builder.Build()))
	server.AddRoute(http.MethodGet, "/user", func(ctx *web.Context) {

	})
}
