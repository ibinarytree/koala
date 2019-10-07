package controller

import (
	"context"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/tools/koala/output/generate/hello"
	"google.golang.org/grpc/metadata"
)

type SayHelloController struct {
}

//检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *SayHelloController) CheckParams(ctx context.Context, r *hello.HelloRequest) (err error) {

	return
}

//SayHello函数的实现
func (s *SayHelloController) Run(ctx context.Context, r *hello.HelloRequest) (
	resp *hello.HelloResponse, err error) {
	resp = &hello.HelloResponse{
		Reply: "hello",
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		logs.AddField(ctx, "client_md", md)
	}

	logs.Debug(ctx, "req=%#v", r)
	logs.AddField(ctx, "use_id", 3838838383)
	return
}
