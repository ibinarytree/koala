package main

var rpcClientTemplate = `
package {{.Package.Name}}c


import (
	"context"
	"fmt"

	"{{.Prefix}}/generate/{{.Package.Name}}"
	
	"github.com/ibinarytree/koala/rpc"
	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/meta"
	
)

type {{Capitalize .Package.Name}}Client struct {
	serviceName string
	client *rpc.KoalaClient
}

func New{{Capitalize .Package.Name}}Client(serviceName string, opts...rpc.RpcOptionFunc) *{{Capitalize .Package.Name}}Client {
	c :=  &{{Capitalize .Package.Name}}Client{
		serviceName: serviceName,
	}
	c.client = rpc.NewKoalaClient(serviceName, opts...)
	return c
}

{{range .Rpc}}
func (s *{{Capitalize $.Package.Name}}Client) {{.Name}}(ctx context.Context, r*{{$.Package.Name}}.{{.RequestType}})(resp*{{$.Package.Name}}.{{.ReturnsType}}, err error){
	/*
	middlewareFunc := rpc.BuildClientMiddleware(mwClient{{.Name}})
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}
*/
	mkResp, err := s.client.Call(ctx, "{{.Name}}", r, mwClient{{.Name}})
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
	/*
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}*/
	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	client := {{$.Package.Name}}.New{{Capitalize $.Package.Name}}ServiceClient(rpcMeta.Conn)
	return client.{{.Name}}(ctx, req)
}
{{end}}
`
