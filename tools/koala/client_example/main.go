package main

//client.go

import (
	"context"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/tools/koala/client_example/generate/client/helloc"
	"github.com/ibinarytree/koala/tools/koala/client_example/generate/hello"
)

const (
	address     = "localhost:8080"
	defaultName = "world"
)

func main() {
	client := helloc.NewHelloClient("hello")
	ctx := context.Background()
	resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
	if err != nil {
		logs.Error(ctx, "could not greet: %v", err)
		return
	}

	logs.Info(ctx, "Greeting: %s", resp.Reply)
	logs.Stop()
}
