package main

var main_template = `
package main
import(
	"log"
	
	"github.com/ibinarytree/koala/server"
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

var routerServer = &router.RouterServer{}
	
func main() {

	err := server.Init("{{.Package.Name}}")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	{{.Package.Name}}.Register{{Capitalize .Package.Name}}ServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
`
