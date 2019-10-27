
package helloc


import (
	"context"
	"fmt"

	"github.com/ibinarytree/koala/tools/koala/client_example/generate/hello"
	
	"github.com/ibinarytree/koala/rpc"
	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/meta"
	
)

type HelloClient struct {
	serviceName string
	client *rpc.KoalaClient
}

func NewHelloClient(serviceName string, opts...rpc.RpcOptionFunc) *HelloClient {
	c :=  &HelloClient{
		serviceName: serviceName,
	}
	c.client = rpc.NewKoalaClient(serviceName, opts...)
	return c
}


func (s *HelloClient) SayHello(ctx context.Context, r*hello.HelloRequest)(resp*hello.HelloResponse, err error){
	/*
	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}
*/
	mkResp, err := s.client.Call(ctx, "SayHello", r, mwClientSayHello)
	if err != nil {
		return nil, err
	}
	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
		return nil, err
	}
	


	return resp, err
}


func mwClientSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	/*
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}*/
	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*hello.HelloRequest)
	client := hello.NewHelloServiceClient(rpcMeta.Conn)
	return client.SayHello(ctx, req)
}

