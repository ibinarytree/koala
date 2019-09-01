package middleware

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow() bool
}

func NewRateLimitMiddleware(l Limiter) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			allow := l.Allow()
			if !allow {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}

			return next(ctx, req)
		}
	}
}
