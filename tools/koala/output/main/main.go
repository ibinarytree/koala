package main

import (
	"context"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/server"

	"github.com/ibinarytree/koala/tools/koala/output/router"

	"github.com/ibinarytree/koala/tools/koala/output/generate/hello"
)

var routerServer = &router.RouterServer{}

func main() {

	err := server.Init("hello")
	if err != nil {
		logs.Error(context.TODO(), "init service failed, err:%v", err)
		logs.Stop()
		return
	}

	hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
