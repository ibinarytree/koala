package middleware

import (
	"context"

	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/loadbalance"
	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/meta"
)

func NewLoadBalanceMiddleware(balancer loadbalance.LoadBalance) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//从ctx获取rpc的metadata
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = errno.NotHaveInstance
				logs.Error(ctx, "not have instance")
				return
			}
			//生成loadbalance的上下文,用来过滤已经选择的节点
			ctx = loadbalance.WithBalanceContext(ctx)
			for {
				rpcMeta.CurNode, err = balancer.Select(ctx, rpcMeta.AllNodes)
				if err != nil {
					return
				}

				logs.Debug(ctx, "select node:%#v", rpcMeta.CurNode)
				rpcMeta.HistoryNodes = append(rpcMeta.HistoryNodes, rpcMeta.CurNode)
				resp, err = next(ctx, req)
				if err != nil {
					//连接错误的话，进行重试
					if errno.IsConnectError(err) {
						continue
					}
					return
				}
				break
			}
			return
		}
	}
}
