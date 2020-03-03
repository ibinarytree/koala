package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ibinarytree/proto"
)

var genMgr *GeneratorMgr = &GeneratorMgr{
	genClientMap: make(map[string]Generator),
	genServerMap: make(map[string]Generator),
	metaData:     &ServiceMetaData{},
}

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf/product",
	"conf/test",
	"model",
	"generate",
	"router",
}

type GeneratorMgr struct {
	genClientMap map[string]Generator
	genServerMap map[string]Generator
	metaData     *ServiceMetaData
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
		proto.WithOption(g.handleOption),
	)

	//必须包含option go_package
	if g.metaData.containGoPackage == false {
		err = fmt.Errorf("must define option go_package=\"xxx\"")
		return
	}

	partLen := len(g.metaData.serviceNameParts)
	if partLen == 0 {
		err = fmt.Errorf("package name is invalid")
		return
	}

	lastPart := g.metaData.serviceNameParts[partLen-1]
	if lastPart != g.metaData.PackageName {
		err = fmt.Errorf("last part of service name must equal package name")
		return
	}

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

func (g *GeneratorMgr) handleOption(r *proto.Option) {

	g.metaData.options = append(g.metaData.options, r)
	if strings.ToLower(r.Name) == "go_package" {
		g.metaData.containGoPackage = true
		g.metaData.PackageName = r.Constant.Source
	}
}

func (g *GeneratorMgr) handlePackage(r *proto.Package) {
	g.metaData.Package = r
	g.metaData.ServiceName = g.metaData.Package.Name
	g.metaData.serviceNameParts = strings.Split(g.metaData.ServiceName, ".")
	g.metaData.ServiceNamePartsPath = strings.Join(g.metaData.serviceNameParts, "/")
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

func (g *GeneratorMgr) initOutputDir(opt *Option) (err error) {

	goPath := os.Getenv("GOPATH")
	if len(opt.Prefix) > 0 {
		//假如用户指定Prefix=github.com/ibinarytree/koala/example
		//outputDir=$GOPATH/$Prefix
		opt.Output = path.Join(goPath, "src", opt.Prefix)
		return
	}

	//如果用户梅有指定包的路径，那么使用当前路径作为包的路径以及output目录
	//exeFilePath = "C:\\xxx\\a.exe"
	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe path:%v", exeFilePath)
		return
	}
	//C:/project/src/xxx/
	opt.Output = strings.ToLower(exeFilePath[0:lastIdx])
	srcPath := strings.ToLower(path.Join(goPath, "src/"))
	if srcPath[len(srcPath)-1] != '/' {
		srcPath = fmt.Sprintf("%s/", srcPath)
	}
	opt.Prefix = strings.Replace(opt.Output, srcPath, "", -1)

	fmt.Printf("opt output:%s, prefix:%s, gopath:%s\n", opt.Output, opt.Prefix, goPath)
	return
}

func (g *GeneratorMgr) fmtCode(opt *Option) {

	cmd := exec.Command("gofmt", "-w", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("warn:gofmt failed, err:%v\n", err)
		return

	}

	//遍历当前目录，把所有go文件都用goimports 格式化
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".go") {
			fmt.Printf("path:%s, name:%s\n", path, info.Name())
			cmd := exec.Command("goimports", "-w", path)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			if err != nil {
				fmt.Printf("warn:goimports failed, err:%v\n", err)
				return err
			}
			return nil
		}

		return nil
	})
	return
}

func (g *GeneratorMgr) Run(opt *Option) (err error) {

	err = g.initOutputDir(opt)
	if err != nil {
		return
	}

	err = g.parseService(opt)
	if err != nil {
		return
	}

	defer func() {
		//gofmt it
		g.fmtCode(opt)
	}()

	g.metaData.Prefix = opt.Prefix
	if opt.GenClientCode {
		for _, gen := range g.genClientMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
		return
	}

	if opt.GenServerCode {
		err = g.createAllDir(opt)
		if err != nil {
			return
		}

		for _, gen := range g.genServerMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
		return
	}

	return
}

func RegisterClientGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genClientMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
		return
	}

	genMgr.genClientMap[name] = gen
	return
}

func RegisterServerGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genServerMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
		return
	}

	genMgr.genServerMap[name] = gen
	return
}
