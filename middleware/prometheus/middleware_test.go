//go:build e2e

package prometheus

import (
	"dengming20240317/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {

	builder := MiddlewareBuilder{
		Namespace: "geekbang_go",
		Subsystem: "web",
		Name:      "http_response",
	}
	server := web.NewHTTPServer(web.ServerWithMiddleWARE(builder.Build()))
	server.AddRoute(http.MethodGet, "/user", func(ctx *web.Context) {

	})
	//监控prometheus的数据
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()

	server.Start(":8081")
}
