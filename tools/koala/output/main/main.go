
package main
import(
	"log"
	
	"github.com/ibinarytree/koala/server"
	
		"github.com/ibinarytree/koala/tools/koala/output/router"
	
	
		"github.com/ibinarytree/koala/tools/koala/output/generate/hello"
	
)

var routerServer = &router.RouterServer{}
	
func main() {

	err := server.Init("hello")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
