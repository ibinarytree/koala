package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type GrpcGenerator struct {
}

func (d *GrpcGenerator) run(opt *Option, metaData *ServiceMetaData, protoFile string) (err error) {

	//protoc --go_out=plugins=grpc:. hello.proto
	dir := path.Join(opt.GoPath, "src")
	os.MkdirAll(dir, 0755)
	outputParams := fmt.Sprintf("plugins=grpc:%s", dir)

	var params []string
	params = append(params, "--go_out")
	params = append(params, outputParams)
	params = append(params, protoFile)
	params = append(params, "-I", path.Dir(protoFile))

	for _, val := range opt.ProtoPaths {
		params = append(params, "-I", val)
	}

	cmd := exec.Command("protoc", params...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("grpc generator failed, err:%v\n", err)
		return

	}

	return
}

func (d *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	err = d.run(opt, metaData, opt.Proto3Filename)
	if err != nil {
		fmt.Printf("generate grpc:%s failed, err:%v\n", opt.Proto3Filename, err)
		return
	}

	for _, file := range opt.ImportFiles {
		err = d.run(opt, metaData, file)
		if err != nil {
			fmt.Printf("generate grpc:%s failed, err:%v\n", file, err)
			return
		}
	}
	return
}

func init() {
	gc := &GrpcGenerator{}

	RegisterServerGenerator("grpc generator", gc)
}
