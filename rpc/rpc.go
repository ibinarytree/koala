package rpc

import (
	"github.com/ibinarytree/koala/middleware"
	_ "github.com/ibinarytree/koala/registry/etcd"
)

func BuildClientMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware
	if len(mids) == 0 {
		return handle
	}

	m := middleware.Chain(mids[0], mids...)
	return m(handle)
}
