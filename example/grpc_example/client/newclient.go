package main

import (
	"context"
	"fmt"

	"github.com/ibinarytree/koala/example/grpc_example/hello"
	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/rpc"
	"google.golang.org/grpc"
)

type HelloClient struct {
	serviceName string
}

func NewHelloClient(serviceName string) *HelloClient {
	return &HelloClient{
		serviceName: serviceName,
	}
}

func (h *HelloClient) SayHelloV1(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(context.Background(), "did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()

	c := hello.NewHelloServiceClient(conn)
	r, err := c.SayHello(ctx, in, opts...)
	if err != nil {
		logs.Error(ctx, "could not greet: %v", err)
		return nil, err
	}
	return r, err
}

func (h *HelloClient) SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {

	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, in)
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}
	req := request.(*hello.HelloRequest)
	defer conn.Close()
	client := hello.NewHelloServiceClient(conn)
	return client.SayHello(ctx, req)
}
