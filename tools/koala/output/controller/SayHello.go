package controller

import (
	"context"
	"fmt"

	"github.com/ibinarytree/koala/logs"

	"github.com/ibinarytree/koala/tools/koala/output/generate/hello"
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
	err = fmt.Errorf(
		"Length of `Name` cannot be more than 10 characters")

	logs.Debug(ctx, "req=%#v", r)
	logs.AddField(ctx, "use_id", 3838838383)
	return
}
