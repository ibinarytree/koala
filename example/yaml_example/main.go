package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Site  SiteConfig  `yaml:"site"`
	Nginx NginxConfig `yaml:"nginx"`
}

type SiteConfig struct {
	Port      int    `yaml:"port"`
	HttpsOn   bool   `yaml:"https_on"`
	Domain    string `yaml:"domain"`
	HttpsPort int    `yaml:"https_port"`
}

type NginxConfig struct {
	Port     int      `yaml:"port"`
	LogPath  string   `yaml:"log_path"`
	SiteName string   `yaml:"site_name"`
	SiteAddr string   `yaml:"site_addr"`
	Upstream []string `yaml:"upstream"`
}

func main() {

	fmt.Printf("os.args[0]=%s\n", os.Args[0])
	data, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Printf("unmarshal failed err:%v\n", err)
		return
	}

	fmt.Printf("site port:%d\n", conf.Site.Port)
	fmt.Printf("conf is %#v\n", conf)
}
