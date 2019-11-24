package main

//client.go

import (
	"context"
	"time"

	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/rpc"
	"github.com/ibinarytree/koala/tools/koala/client_example/generate/client/helloc"
	"github.com/ibinarytree/koala/tools/koala/client_example/generate/hello"
)

const (
	address     = "localhost:8080"
	defaultName = "world"
)

func main() {
	client := helloc.NewHelloClient("hello", rpc.WithLimitQPS(5),
		rpc.WithClientServiceName("hello-client-example"))
	var count int
	for {
		count++
		ctx := context.Background()
		resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
		if err != nil {
			if count%100 == 0 {
				logs.Error(ctx, "could not greet: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if count%100 == 0 {
			logs.Info(ctx, "Greeting: %s", resp.Reply)
		}

		time.Sleep(100 * time.Millisecond)
	}
	logs.Stop()
}
