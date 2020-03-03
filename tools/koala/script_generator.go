package main

import (
	"fmt"
	"os"
	"path"
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
	err = d.RendScript(opt, metaData, filename, build_template)
	if err != nil {
		fmt.Printf("render:%s failed, err:%v\n", filename, err)
		return
	}

	filename = "scripts/start.sh"
	err = d.RendScript(opt, metaData, filename, start_template)
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
