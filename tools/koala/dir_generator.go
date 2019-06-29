package main

import (
	"fmt"
	"os"
	"path"
)

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
}

type DirGenerator struct {
	dirList []string
}

func (d *DirGenerator) Run(opt *Option) (err error) {

	for _, dir := range d.dirList {
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("mkdir dir %s failed, err:%v\n", dir, err)
			return
		}
	}
	return
}

func init() {
	dir := &DirGenerator{
		dirList: AllDirList,
	}

	Register("dir generator", dir)
}
