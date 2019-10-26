package meta

import (
	"context"

	"github.com/ibinarytree/koala/registry"
)

type RpcMeta struct {
	//调用方名字
	Caller string
	//服务提供方
	ServiceName string
	//调用的方法
	Method string
	//调用方集群
	CallerCluster string
	//服务提供方集群
	ServiceCluster string
	//TraceID
	TraceID string
	//环境
	Env string
	//调用方IDC
	CallerIDC string
	//服务提供方IDC
	ServiceIDC string
	//当前节点
	CurNode *registry.Node
	//历史选择节点
	HistoryNodes []*registry.Node
	//服务提供方的节点列表
	AllNodes []*registry.Node
	//
}

type rpcMetaContextKey struct{}

func GetRpcMeta(ctx context.Context) *RpcMeta {
	meta, ok := ctx.Value(rpcMetaContextKey{}).(*RpcMeta)
	if !ok {
		meta = &RpcMeta{}
	}

	return meta
}

func InitRpcMeta(ctx context.Context, service, method, caller string) context.Context {
	meta := &RpcMeta{
		Method:      method,
		ServiceName: service,
		Caller:      caller,
	}
	return context.WithValue(ctx, rpcMetaContextKey{}, meta)
}
