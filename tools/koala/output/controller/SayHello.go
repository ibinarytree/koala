
package controller
import(
	"context"
	
		"github.com\ibinarytree\koala\tools\koala\output/hello"
	
)

type SayHelloController struct {
}


//检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *SayHelloController) CheckParams(ctx context.Context, r*hello.HelloRequest) (err error) {
	return
}

//SayHello函数的实现
func (s *SayHelloController) Run(ctx context.Context, r*hello.HelloRequest) (
	resp*hello.HelloResponse, err error) {
	return
}

