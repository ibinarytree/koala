package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

type RpcClientGenerator struct {
}

func (d *RpcClientGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	packageName := fmt.Sprintf("%sc", metaData.Package.Name)
	dir := path.Join(opt.Output, "generate", "client", packageName)
	os.MkdirAll(dir, 0755)

	filename := path.Join(dir, "client.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}

	defer file.Close()
	err = d.render(file, rpcClientTemplate, metaData)
	if err != nil {
		fmt.Printf("render failed, err:%v\n", err)
		return
	}

	return
}

func (d *RpcClientGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {

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
