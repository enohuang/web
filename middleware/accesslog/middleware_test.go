package accesslog

import (
	"dengming20240317/web"
	"fmt"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_LogFunc(t *testing.T) {

	build := MiddlewareBuilder{}
	mdl := build.LogFunc(func(log string) {
		fmt.Println("ccccccccccccccccccccccc")
		fmt.Println(log)
	}).Build()

	server := web.NewHTTPServer(web.ServerWithMiddleWARE(mdl))
	server.Mdls()
	server.AddRoute(http.MethodGet, "/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it is me")
	})
	req, err := http.NewRequest(http.MethodGet, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
