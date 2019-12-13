package main

import (
	"log"

	"github.com/ibinarytree/koala/server"

	"github.com/ibinarytree/koala/tools/koala/output/router"

	"github.com/ibinarytree/koala/tools/koala/output/generate/com/google/hello"
)

var routerServer = &router.RouterServer{}

func main() {

	err := server.Init("com.google.hello")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
