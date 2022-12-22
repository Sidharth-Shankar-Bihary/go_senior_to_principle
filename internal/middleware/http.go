package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/opentracing/opentracing-go"
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/utils"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

// InjectTracer needs a param which is *env.Environment type.
// InjectTracer return is also a func, which named func(http.Handler) and the type is http.Handler.
func InjectTracer(e *env.Environment) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					e.L().Error(
						"!!!!!!!! panic recover !!!!!!!!",
						zap.String("err", fmt.Sprintf("%+v", err)),
						zap.ByteString("trace", debug.Stack()),
					)
				}
			}()

			tracer := opentracing.GlobalTracer()
			// Create a span A
			parentSpan := tracer.StartSpan(r.URL.Path)

			// Get TraceID
			if sc, ok := parentSpan.Context().(jaeger.SpanContext); ok {
				e.Log.With(zap.Any("trace-id", sc.TraceID()))
				// fmt.Println("trace-id", sc.TraceID())
			}

			defer parentSpan.Finish()
			ctx := context.WithValue(r.Context(), utils.SpanStr, parentSpan)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// multiple spans not work yet
// func WrapTrace(f gin.HandlerFunc, funcName string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		span := utils.StartSpan(c, funcName)
// 		if span != nil {
// 			defer span.Finish()
// 		}
// 		ctx := context.WithValue(c, utils.SpanStr, span)
// 		c.Request.WithContext(ctx)
// 		f(c)
// 	}
// }
