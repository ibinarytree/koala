package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ibinarytree/koala/util"
	"github.com/ibinarytree/proto"
)

type CtrlGenerator struct {
}

type RpcMeta struct {
	Rpc *proto.RPC
	//Package *proto.Package
	//Prefix  string
	*ServiceMetaData
}

func (d *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", opt.Proto3Filename, err)
		return
	}

	defer reader.Close()
	return d.generateRpc(opt, metaData)
}

func (d *CtrlGenerator) generateRpc(opt *Option, metaData *ServiceMetaData) (err error) {

	for _, rpc := range metaData.Rpc {
		var file *os.File
		tmpName := ToUnderScoreString(rpc.Name)
		filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", strings.ToLower(tmpName)))
		fmt.Printf("filename is %s\n", filename)
		exist := util.IsFileExist(filename)
		if exist {
			continue
		}
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("open file:%s failed, err:%v\n", filename, err)
			return
		}

		rpcMeta := &RpcMeta{}
		rpcMeta.Rpc = rpc
		rpcMeta.ServiceMetaData = metaData
		//rpcMeta.Package = metaData.Package
		//rpcMeta.Prefix = metaData.Prefix

		err = d.render(file, controller_template, rpcMeta)
		if err != nil {
			fmt.Printf("render controller failed err:%v\n", err)
			return
		}
		defer file.Close()
	}
	return
}

func (d *CtrlGenerator) render(file *os.File, data string, metaData *RpcMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	ctrl := &CtrlGenerator{}
	RegisterServerGenerator("ctrl generator", ctrl)
}
