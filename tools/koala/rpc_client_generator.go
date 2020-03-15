package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

type RpcClientGenerator struct {
}

type RpcClientData struct {
	*ServiceMetaData
	ClientImportPath  string
	ClientPackageName string
}

func (d *RpcClientGenerator) run(opt *Option, rpcClientData *RpcClientData, templateFile, filename string) (err error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}

	defer file.Close()
	err = d.render(file, templateFile, rpcClientData)
	if err != nil {
		fmt.Printf("render failed, err:%v\n", err)
		return
	}

	return
}

func (d *RpcClientGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	//所有生成的rpc client的包名都以c后缀结尾，最后一级目录加上c的后缀
	packagePath := metaData.ImportPath
	if packagePath[len(packagePath)-1] == '/' || packagePath[len(packagePath)-1] == '\\' {
		packagePath = packagePath[:len(packagePath)-1]
	}

	packagePath = fmt.Sprintf("%sc", packagePath)
	dir := path.Join(opt.Output, "rpc/krpc/clients", packagePath)
	os.MkdirAll(dir, 0755)

	rpcClientData := &RpcClientData{
		ServiceMetaData:   metaData,
		ClientImportPath:  path.Join(metaData.Prefix, "rpc/krpc/clients", packagePath),
		ClientPackageName: fmt.Sprintf("%sc", metaData.PackageName),
	}

	//1. generate koala client
	filename := path.Join(dir, "client.go")
	err = d.run(opt, rpcClientData, rpcClientTemplate, filename)
	if err != nil {
		fmt.Printf("generate clients failed, filename:%s, err:%v\n", filename, err)
		return
	}

	//2. generate koala client wrap
	filename = path.Join(opt.Output, "rpc/krpc/", fmt.Sprintf("%s_client_wrap.go", metaData.PackageName))
	err = d.run(opt, rpcClientData, grpcClientWrapTemplate, filename)
	if err != nil {
		fmt.Printf("generate clients  wrap failed, filename:%s, err:%v\n", filename, err)
		return
	}

	return
}

func (d *RpcClientGenerator) render(file *os.File, data string, metaData *RpcClientData) (err error) {

	t := template.New("main").Funcs(templateFuncMap)
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	rpcClient := &RpcClientGenerator{}
	RegisterClientGenerator("rpc client generator", rpcClient)
}
