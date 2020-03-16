package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	var opt Option
	var importFiles cli.StringSlice
	var protoPaths cli.StringSlice

	app := cli.NewApp()
	app.Version = "2.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "idl filename, May be specified multiple times",
			Required:    true,
			Destination: &opt.Proto3Filename,
		},
		cli.BoolFlag{
			Name:        "c",
			Usage:       "generate grpc client code",
			Destination: &opt.GenClientCode,
		},
		cli.BoolFlag{
			Name:        "s",
			Usage:       "generate grpc server code",
			Destination: &opt.GenServerCode,
		},
		cli.StringFlag{
			Name:        "p",
			Value:       "",
			Usage:       "prefix of package",
			Destination: &opt.Prefix,
		},
		cli.StringSliceFlag{
			Name:  "i",
			Usage: "import proto file, Specify the proto file in which for proto file imports.May be specified multiple times",
			Value: &importFiles,
		},
		cli.StringSliceFlag{
			Name:  "proto_path",
			Usage: "Specify the directory in which to search for imports.  May be specified multiple times",
			Value: &protoPaths,
		},
	}

	app.Action = func(c *cli.Context) error {

		for _, file := range importFiles {
			opt.ImportFiles = append(opt.ImportFiles, file)
		}
		for _, proto := range protoPaths {
			opt.ProtoPaths = append(opt.ProtoPaths, proto)
		}

		//命令行程序代码的入口
		err := genMgr.Run(&opt)
		if err != nil {
			fmt.Printf("code generator failed, err:%v\n", err)
			return err
		}

		fmt.Println("code generate succ")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
