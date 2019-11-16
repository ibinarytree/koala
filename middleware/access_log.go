package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/meta"
	"google.golang.org/grpc/status"
)

func AccessLogMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		startTime := time.Now()
		resp, err = next(ctx, req)

		serverMeta := meta.GetServerMeta(ctx)
		errStatus, _ := status.FromError(err)

		cost := time.Since(startTime).Nanoseconds() / 1000
		logs.AddField(ctx, "cost_us", cost)
		logs.AddField(ctx, "method", serverMeta.Method)

		logs.AddField(ctx, "cluster", serverMeta.Cluster)
		logs.AddField(ctx, "env", serverMeta.Env)
		logs.AddField(ctx, "server_ip", serverMeta.ServerIP)
		logs.AddField(ctx, "client_ip", serverMeta.ClientIP)
		logs.AddField(ctx, "idc", serverMeta.IDC)
		logs.Access(ctx, "result=%v", errStatus.Code())

		return
	}
}

func RpcLogMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {

		ctx = logs.WithFieldContext(ctx)
		startTime := time.Now()
		resp, err = next(ctx, req)

		rpcMeta := meta.GetRpcMeta(ctx)
		errStatus, _ := status.FromError(err)

		cost := time.Since(startTime).Nanoseconds() / 1000
		logs.AddField(ctx, "cost_us", cost)
		logs.AddField(ctx, "method", rpcMeta.Method)
		logs.AddField(ctx, "server", rpcMeta.ServiceName)

		logs.AddField(ctx, "caller_cluster", rpcMeta.CallerCluster)
		logs.AddField(ctx, "upstream_cluster", rpcMeta.ServiceCluster)
		logs.AddField(ctx, "rpc", 1)
		logs.AddField(ctx, "env", rpcMeta.Env)

		var upstreamInfo string
		for _, node := range rpcMeta.HistoryNodes {
			upstreamInfo += fmt.Sprintf("%s:%d,", node.IP, node.Port)
		}

		logs.AddField(ctx, "upstream", upstreamInfo)
		logs.AddField(ctx, "caller_idc", rpcMeta.CallerIDC)
		logs.AddField(ctx, "upstream_idc", rpcMeta.ServiceIDC)
		logs.Access(ctx, "result=%v", errStatus.Code())

		return
	}
}
