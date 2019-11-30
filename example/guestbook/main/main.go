
package main
import(
	"log"
	
	"github.com/ibinarytree/koala/server"
	
		"github.com/ibinarytree/koala/example/guestbook/router"
	
	
		"github.com/ibinarytree/koala/example/guestbook/generate/guestbook"
	
)

var routerServer = &router.RouterServer{}
	
func main() {

	err := server.Init("guestbook")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	guestbook.RegisterGuestBookServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
