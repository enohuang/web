//go:build e2e

package opentelemetry

import (
	"dengming20240317/web"
	"go.opentelemetry.io/otel"
	"net/http"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	builder := MiddlewareBuilder{
		Tracer: tracer,
	}

	server := web.NewHTTPServer(web.ServerWithMiddleWARE(builder.Build()))

	server.AddRoute(http.MethodGet, "/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()

		c, third2 := tracer.Start(c, "third_layer_2  ")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()

		c, first := tracer.Start(ctx.Req.Context(), "first_layer_1")
		defer first.End()
		ctx.Resp.Write([]byte("hello   "))

	})

	server.Start(":8081")
}
