package logger

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

//type Middleware func(Endpoint) Endpoint

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("msg", "calling endpoint")
			start := time.Now()
			defer func() {
				logger.Log("msg", "called endpoint", "duration", time.Since(start))
			}()
			return next(ctx, request)
		}
	}
}
