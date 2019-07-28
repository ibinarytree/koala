
package router
import(
	"context"
	
		"github.com\ibinarytree\koala\tools\koala\output/hello"
	
)
	
type RouterServer struct{}


func (s *RouterServer) SayHello(ctx context.Context, r*hello.HelloRequest)(resp*hello.HelloResponse){
	/*inst := &SayHelloController{}
	err = inst.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = inst.Run(ctx, r)*/
	return
}

