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
		"generate/{{.ServiceNamePartsPath}}"
	{{else}}
		"{{.Prefix}}/generate/{{.ServiceNamePartsPath}}"
	{{end}}
)

var routerServer = &router.RouterServer{}
	
func main() {

	err := server.Init("{{.ServiceName}}")
	if err != nil {
		log.Fatal("init service failed, err:%v", err)
		return
	}

	{{.PackageName}}.Register{{Capitalize .PackageName}}ServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}
`
