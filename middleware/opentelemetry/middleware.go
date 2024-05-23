package opentelemetry

import (
	"dengming20240317/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "dengming20240317/web/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

/*func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		Tracer: tracer,
	}
}*/

func (m *MiddlewareBuilder) Build() web.Middleware {

	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}

	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			//尝试与客户端的tarce 结合在一起

			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			_, sqan := m.Tracer.Start(reqCtx, "unknown" /*命中的路由*/)
			defer sqan.End()

			sqan.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			sqan.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			sqan.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			sqan.SetAttributes(attribute.String("http.host", ctx.Req.URL.Host))
			//直接调用下一步
			ctx.Req = ctx.Req.WithContext(reqCtx)
			next(ctx)

			sqan.SetName(ctx.MatchedRoute)
			sqan.SetAttributes(attribute.Int("http.status", 200))
		}
	}
}
