package main

var router_template = `
package router
import(
	"context"

	"github.com/ibinarytree/koala/server"
	"github.com/ibinarytree/koala/meta"
	{{if not .Prefix}}
		"generate/{{.ServiceNamePartsPath}}"
	{{else}}
		"{{.Prefix}}/generate/{{.ServiceNamePartsPath}}"
	{{end}}
	{{if not .Prefix}}
	"controller"
{{else}}
	"{{.Prefix}}/controller"
{{end}}
)
	
type RouterServer struct{}

{{range .Rpc}}
func (s *RouterServer) {{.Name}}(ctx context.Context, r*{{$.PackageName}}.{{.RequestType}})(resp*{{$.PackageName}}.{{.ReturnsType}}, err error){
	
	ctx = meta.InitServerMeta(ctx,"{{$.PackageName}}", "{{.Name}}")
	mwFunc := server.BuildServerMiddleware(mw{{.Name}})
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	
	resp = mwResp.(*{{$.PackageName}}.{{.ReturnsType}})
	return
}


func mw{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	
		r := request.(*{{$.PackageName}}.{{.RequestType}})
		ctrl := &controller.{{.Name}}Controller{}
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
}
{{end}}
`
