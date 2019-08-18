package main

//client.go

import (
	"log"
	"os"
	"time"

	pb "github.com/ibinarytree/koala/example/grpc_example/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:12345"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	for {
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatal("could not greet: %v", err)
		}
		_ = r
		//log.Printf("Greeting: %s", r.Reply)
		time.Sleep(time.Millisecond * 10)
	}
}
