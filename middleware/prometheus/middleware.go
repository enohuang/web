package prometheus

import (
	"dengming20240317/web"
	"strconv"
	"time"
)
import "github.com/prometheus/client_golang/prometheus"

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

// func(next HandleFunc) HandleFunc
func (m *MiddlewareBuilder) Build() web.Middleware {

	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Namespace: m.Namespace,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(vector)

	return func(next web.HandleFunc) web.HandleFunc {

		return func(ctx *web.Context) {

			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()
				pattern := ctx.MatchedRoute
				if pattern == "" {
					pattern = "unknown"
				}
				//status.Code
				vector.WithLabelValues(pattern, ctx.Req.Method, strconv.Itoa(1)).Observe(float64(duration))
			}()
			next(ctx)

		}
	}
}
