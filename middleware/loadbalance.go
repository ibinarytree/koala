package middleware

import (
	"context"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/meta"
	"github.com/ibinarytree/koala/loadbalance"
)

func NewLoadBalanceMiddleware(balancer loadbalance.LoadBalance) Middleware{
return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//从ctx获取rpc的metadata
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = errno.NotHaveInstance
				logs.Error(ctx, "not have instance")
				return
			}

			for {
				resp, err = next(ctx, req)
			}
			return
		}
}
}
