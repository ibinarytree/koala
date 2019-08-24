package main

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/ibinarytree/koala/util"
)

type ConfigGenerator struct {
}

func (d *ConfigGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	err = d.generateByEnv(util.PRODUCT_ENV, opt, metaData)
	if err != nil {
		fmt.Printf("generate config failed, err:%v\n", err)
		return
	}

	err = d.generateByEnv(util.TEST_ENV, opt, metaData)
	if err != nil {
		fmt.Printf("generate config failed, err:%v\n", err)
		return
	}

	return
}

func (d *ConfigGenerator) generateByEnv(env string, opt *Option, metaData *ServiceMetaData) (err error) {

	var file *os.File
	filename := path.Join(opt.Output, "conf", env, fmt.Sprintf("%s.yaml", metaData.Package.Name))
	fmt.Printf("filename is %s\n", filename)

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}

	err = d.render(file, config_template, metaData)
	if err != nil {
		fmt.Printf("render controller failed err:%v\n", err)
		return
	}
	defer file.Close()
	return
}

func (d *ConfigGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	conf := &ConfigGenerator{}
	Register("conf generator", conf)
}
