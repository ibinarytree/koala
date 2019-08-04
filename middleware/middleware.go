package middleware

import (
	"context"
)

type MiddlewareFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

type Middleware func(MiddlewareFunc) MiddlewareFunc

var userMiddleware []Middleware

// Chain is a helper function for composing middlewares. Requests will
// traverse them in the order they're declared. That is, the first middleware
// is treated as the outermost middleware.
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func Use(m ...Middleware) {
	userMiddleware = append(userMiddleware, m...)
}

func BuildServerMiddleware(handle MiddlewareFunc) (handleChain MiddlewareFunc) {

	var mids []Middleware
	if len(userMiddleware) != 0 {
		mids = append(mids, userMiddleware...)
	}

	if len(mids) > 0 {
		m := Chain(mids[0], mids[1:]...)
		return m(handle)
	}

	return handle
}
