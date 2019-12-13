package router

import (
	"context"

	"github.com/ibinarytree/koala/meta"
	"github.com/ibinarytree/koala/server"

	"github.com/ibinarytree/koala/tools/koala/output/generate/com/google/hello"

	"github.com/ibinarytree/koala/tools/koala/output/controller"
)

type RouterServer struct{}

func (s *RouterServer) SayHello(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error) {

	ctx = meta.InitServerMeta(ctx, "hello", "SayHello")
	mwFunc := server.BuildServerMiddleware(mwSayHello)
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}

	resp = mwResp.(*hello.HelloResponse)
	return
}

func mwSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {

	r := request.(*hello.HelloRequest)
	ctrl := &controller.SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}
