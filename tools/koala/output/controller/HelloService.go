package controller
import(
"context"
hello "github.com/ibinarytree/koala/tools/koala/output/generate"
)

type Server struct{}


func (s *Server) SayHello(ctx context.Context, r*hello.HelloRequest)(resp*hello.HelloResponse, err error){
return
}

