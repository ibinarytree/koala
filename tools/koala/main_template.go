package main

var main_template = `
package main
import(
	"net"
	"log"
	"google.golang.org/grpc"
	
	{{if not .Prefix}}
	"router"
	{{else}}
		"{{.Prefix}}/router"
	{{end}}

	
	{{if not .Prefix}}
		"generate/{{.Package.Name}}"
	{{else}}
		"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
)

var server = &router.RouterServer{}
var port= ":12345"
	
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.Register{{.Service.Name}}Server(s, server)
	s.Serve(lis)
}
`
