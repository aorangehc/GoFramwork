package ratelimit

import (
	"context"

	"github.com/baidubce/bce-sdk-go/rate"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RateLimitMiddleware(limiter *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if !limiter.Allow() {
				return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
			}
			return next(ctx, request)
		}
	}
}
