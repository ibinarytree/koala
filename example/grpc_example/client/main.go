package main

//client.go

import (
	"context"
	"os"
	"time"

	"github.com/ibinarytree/koala/example/grpc_example/hello"
	//"golang.org/x/net/context"
	"github.com/ibinarytree/koala/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:8080"
	defaultName = "world"
)

func rawClientExample() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(context.Background(), "did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := hello.NewHelloServiceClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	for {
		ctx := context.Background()
		ctx = metadata.AppendToOutgoingContext(ctx, "koala_trace_id", "888888888888888888888888888")
		r, err := c.SayHello(ctx, &hello.HelloRequest{Name: name})
		if err != nil {
			logs.Error(ctx, "could not greet: %v", err)
			continue
		}
		_ = r
		logs.Error(ctx, "Greeting: %s", r.Reply)
		time.Sleep(time.Millisecond * 10)
	}
}

func myClientExample() {
	client := NewHelloClient("hello")
	ctx := context.Background()
	//rpc client.SayHelloV1的第一版封装
	//resp, err := client.SayHelloV1(ctx, &hello.HelloRequest{Name: "test my client"})
	//client.SayHello基于中间件架构的封装
	resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
	if err != nil {
		logs.Error(ctx, "could not greet: %v", err)
		return
	}

	logs.Info(ctx, "Greeting: %s", resp.Reply)
	return
}

func main() {
	//使用grpc原生client进行测试
	//rawClientExample()
	//使用我们封装的client进行测试
	myClientExample()
	logs.Stop()
}
