package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"text/template"
)

type ScriptGenerator struct {
}

func (d *ScriptGenerator) RendScript(opt *Option, metaData *ServiceMetaData, output, templateData string) (err error) {

	filename := path.Join(opt.Output, output)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}

	defer file.Close()
	err = d.render(file, templateData, metaData)
	if err != nil {
		fmt.Printf("render failed, err:%v\n", err)
		return
	}
	return
}

func (d *ScriptGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {


	filename := "build.sh"
	buildTemplate := build_template
	if runtime.GOOS == "windows" {
		filename = "build.bat"
		buildTemplate = window_build_template
	}

	err = d.RendScript(opt, metaData, filename, buildTemplate)
	if err != nil {
		fmt.Printf("render:%s failed, err:%v\n", filename, err)
		return
	}

	filename = "scripts/start.sh"
	buildTemplate = start_template
	if runtime.GOOS == "windows" {
		filename = "scripts/start.bat"
		buildTemplate = window_start_template
	}

	err = d.RendScript(opt, metaData, filename, buildTemplate)
	if err != nil {
		fmt.Printf("render:%s failed, err:%v\n", filename, err)
		return
	}

	return
}

func (d *ScriptGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	gen := &ScriptGenerator{}
	RegisterServerGenerator("ScriptGenerator", gen)
}
