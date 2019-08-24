package main

var main_template = `
package main
import(
	"net"
	"log"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	if server.GetConf().Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", server.GetConf().Prometheus.Port)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d",server.GetConf().Port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.Register{{.Service.Name}}Server(s, routerServer)
	s.Serve(lis)
}
`
