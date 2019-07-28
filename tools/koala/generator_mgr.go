package main

import (
	"fmt"
	"os"
	"path"

	"github.com/ibinarytree/proto"
)

var genMgr *GeneratorMgr = &GeneratorMgr{
	genMap:   make(map[string]Generator),
	metaData: &ServiceMetaData{},
}

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf",
	"app/router",
	"app/config",
	"model",
	"generate",
	"router",
}

type GeneratorMgr struct {
	genMap   map[string]Generator
	metaData *ServiceMetaData
}

func (g *GeneratorMgr) parseService(opt *Option) (err error) {

	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", opt.Proto3Filename, err)
		return
	}

	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Printf("parse file:%s failed, err:%v\n", opt.Proto3Filename, err)
		return
	}

	proto.Walk(definition,
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage),
	)

	return
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	g.metaData.Messages = append(g.metaData.Messages, m)
}

func (g *GeneratorMgr) handleRPC(r *proto.RPC) {
	g.metaData.Rpc = append(g.metaData.Rpc, r)
}

func (g *GeneratorMgr) handlePackage(r *proto.Package) {
	g.metaData.Package = r
}

func (g *GeneratorMgr) createAllDir(opt *Option) (err error) {

	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("mkdir dir %s failed, err:%v\n", dir, err)
			return
		}
	}
	return
}

func (g *GeneratorMgr) Run(opt *Option) (err error) {

	err = g.parseService(opt)
	if err != nil {
		return
	}

	err = g.createAllDir(opt)
	if err != nil {
		return
	}

	g.metaData.Prefix = opt.Prefix
	for _, gen := range g.genMap {
		err = gen.Run(opt, g.metaData)
		if err != nil {
			return
		}
	}
	return
}

func Register(name string, gen Generator) (err error) {
	_, ok := genMgr.genMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
		return
	}

	genMgr.genMap[name] = gen
	return
}
