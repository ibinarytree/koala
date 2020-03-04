package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type GrpcGenerator struct {
}

func (d *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	//protoc --go_out=plugins=grpc:. hello.proto
	dir := path.Join(opt.Output, "generate", path.Join(metaData.serviceNameParts...))
	os.MkdirAll(dir, 0755)
	outputParams := fmt.Sprintf("plugins=grpc:%s", dir)

	idlDir, idlFilename := filepath.Split(opt.Proto3Filename)

	cmd := exec.Command("protoc", "--go_out", outputParams, "--proto_path", idlDir, idlFilename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("grpc generator failed, err:%v\n", err)
		return

	}
	return
}

func init() {
	gc := &GrpcGenerator{}

	RegisterServerGenerator("grpc generator", gc)
}
