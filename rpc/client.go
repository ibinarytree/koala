package rpc

import (
	"context"
	"time"

	"github.com/ibinarytree/koala/meta"
	"github.com/ibinarytree/koala/middleware"
)

type KoalaClient struct {
	opts *RpcOptions
}

func NewKoalaClient(optfunc ...RpcOptionFunc) *KoalaClient {
	client := &KoalaClient{
		opts: &RpcOptions{
			ConnTimeout:  DefaultConnTimeout,
			WriteTimeout: DefaultWriteTimeout,
			ReadTimeout:  DefaultReadTimeout,
		},
	}

	for _, opt := range optfunc {
		opt(client.opts)
	}

	return client
}

func (k *KoalaClient) getCaller(ctx context.Context) string {

	serverMeta := meta.GetServerMeta(ctx)
	if serverMeta == nil {
		return ""
	}
	return serverMeta.ServiceName
}

func (k *KoalaClient) buildMiddleware(handle middlewareFunc) MiddlewareFunc {

	var mids []middleware.Middleware
	if len(mids) == 0 {
		return handle
	}

	m := middleware.Chain(mids[0], mids...)
	return m(handle)

}

func (k *KoalaClient) Call(ctx context.Context, method string, r interface{}, handle middleware.MiddlewareFunc) (resp interface{}, err error) {

	//构建中间件
	caller := k.getCaller(ctx)
	ctx = meta.InitRpcMeta(ctx, k.opts.ServiceName, method, caller)
	middlewareFunc := k.buildMiddleware(handle)
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp, err
}
