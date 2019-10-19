package main

var rpcClientTemplate = `
package {{.Package.Name}}c


import (
	"context"
	"fmt"

	"{{.Prefix}}/generate/{{.Package.Name}}"
	"github.com/ibinarytree/koala/logs"
	"github.com/ibinarytree/koala/rpc"
	"google.golang.org/grpc"
)

type {{Capitalize .Package.Name}}Client struct {
	serviceName string
}

func New{{Capitalize .Package.Name}}Client(serviceName string) *{{Capitalize .Package.Name}}Client {
	return &{{Capitalize .Package.Name}}Client{
		serviceName: serviceName,
	}
}

{{range .Rpc}}
func (s *{{Capitalize $.Package.Name}}Client) {{.Name}}(ctx context.Context, r*{{$.Package.Name}}.{{.RequestType}})(resp*{{$.Package.Name}}.{{.ReturnsType}}, err error){
	
	middlewareFunc := rpc.BuildClientMiddleware(mwClient{{.Name}})
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*{{$.Package.Name}}.{{.ReturnsType}})
	if !ok {
		err = fmt.Errorf("invalid resp, not *{{$.Package.Name}}.{{.ReturnsType}}")
		return nil, err
	}

	return resp, err
}


func mwClient{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}

	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	defer conn.Close()
	client := {{$.Package.Name}}.New{{Capitalize $.Package.Name}}ServiceClient(conn)
	return client.{{.Name}}(ctx, req)
}
{{end}}
`
