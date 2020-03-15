package main

var controller_template = `
package controller
import(
	"context"
	"{{.ImportPath}}"
)

type {{.Rpc.Name}}Controller struct {
}


//检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *{{.Rpc.Name}}Controller) CheckParams(ctx context.Context, r*{{.PackageName}}.{{.Rpc.RequestType}}) (err error) {
	return
}

//SayHello函数的实现
func (s *{{.Rpc.Name}}Controller) Run(ctx context.Context, r*{{.PackageName}}.{{.Rpc.RequestType}}) (
	resp*{{.PackageName}}.{{.Rpc.ReturnsType}}, err error) {
	return
}

`
