
package helloc


import (
	"context"
	"fmt"

	"github.com/ibinarytree/koala/tools/koala/client_example/generate/hello"
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


func (s *HelloClient) SayHello(ctx context.Context, r*hello.HelloRequest)(resp*hello.HelloResponse, err error){
	
	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, r)
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
	
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}

	req := request.(*hello.HelloRequest)
	defer conn.Close()
	client := hello.NewHelloServiceClient(conn)
	return client.SayHello(ctx, req)
}

