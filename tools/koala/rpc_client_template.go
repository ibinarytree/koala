package main

var rpcClientTemplate = `
package {{.PackageName}}c


import (
	"context"
	"fmt"

	"{{.ImportPath}}"
	"github.com/ibinarytree/koala/rpc"
	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/meta"
	
)

type {{Capitalize .PackageName}}Client struct {
	serviceName string
	client *rpc.KoalaClient
}

func New{{Capitalize .PackageName}}Client(serviceName string, opts...rpc.RpcOptionFunc) *{{Capitalize .PackageName}}Client {
	c :=  &{{Capitalize .PackageName}}Client{
		serviceName: serviceName,
	}
	c.client = rpc.NewKoalaClient(serviceName, opts...)
	return c
}

{{range .Rpc}}
func (s *{{Capitalize $.PackageName}}Client) {{.Name}}(ctx context.Context, r*{{$.PackageName}}.{{.RequestType}})(resp*{{$.PackageName}}.{{.ReturnsType}}, err error){
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
	resp, ok := mkResp.(*{{$.PackageName}}.{{.ReturnsType}})
	if !ok {
		err = fmt.Errorf("invalid resp, not *{{$.PackageName}}.{{.ReturnsType}}")
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

	req := request.(*{{$.PackageName}}.{{.RequestType}})
	client := {{$.PackageName}}.New{{$.Service.Name}}Client(rpcMeta.Conn)

	return client.{{.Name}}(ctx, req)
}
{{end}}
`

var grpcClientWrapTemplate = `
package krpc

import (
	"context"
	"fmt"

	"{{.ClientImportPath}}"
	"{{.ImportPath}}"
)

var (
	k{{Capitalize .ClientPackageName}}Client *{{.ClientPackageName}}.{{.Capitalize .PackageName}}Client
)

func init() {

	k{{Capitalize .ClientPackageName}}Client = {{.ClientPackageName}}.New{{Capitalize .PackageName}}Client("{{.ServiceName}}")
}

{{range .Rpc}}
func {{Capitalize $.PackageName}}{{.Name}}(ctx context.Context, r*{{$.PackageName}}.{{.RequestType}})(resp*{{$.PackageName}}.{{.ReturnsType}}, err error){
	return k{{Capitalize $.ClientPackageName}}Client.{{.Name}}(ctx, r)
}
{{end}}
`
