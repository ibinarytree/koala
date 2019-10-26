package middleware

import (
	"context"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/meta"
	"github.com/ibinarytree/koala/registry"
)

func NewDiscoveryMiddleware(discovery registry.Registry) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//从ctx获取rpc的metadata
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) > 0 {
				return next(ctx, req)
			}

			service, err := discovery.GetService(ctx, rpcMeta.ServiceName)
			if err != nil {
				logs.Error(ctx, "discovery service:%s failed, err:%v", rpcMeta.ServiceName, err)
				return
			}

			rpcMeta.AllNodes = service.Nodes
			resp, err = next(ctx, req)
			return
		}
	}
}
