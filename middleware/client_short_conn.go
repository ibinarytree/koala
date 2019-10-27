package middleware

import (
	"context"
	"fmt"

	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/meta"
	"google.golang.org/grpc"
)

func ShortConnectMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		//从ctx获取rpc的metadata
		rpcMeta := meta.GetRpcMeta(ctx)
		if rpcMeta.CurNode == nil{
			err = errno.InvalidNode
			logs.Error(ctx, "invalid instance")
			return
		}

		address := fmt.Sprintf("%s:%d", rpcMeta.CurNode.IP, rpcMeta.CurNode.Port)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			logs.Error(ctx, "connect %s failed, err:%v", address, err)
			return nil, errno.ConnFailed
		}

		rpcMeta.Conn = conn
		defer conn.Close()
		resp, err = next(ctx, req)
		return
	}
}
