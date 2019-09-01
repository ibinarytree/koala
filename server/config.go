package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ibinarytree/koala/util"
	"gopkg.in/yaml.v2"
)

var (
	koalaConf = &KoalaConf{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "koala_server",
		Regiser: RegisterConf{
			SwitchOn: false,
		},
		Log: LogConf{
			Level: "debug",
			Dir:   "./logs/",
		},
		Limit: LimitConf{
			SwitchOn: true,
			QPSLimit: 50000,
		},
	}
)

type KoalaConf struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	ServiceName string         `yaml:"service_name"`
	Regiser     RegisterConf   `yaml:"register"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`

	//内部的配置项
	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type LimitConf struct {
	QPSLimit int  `yaml:"qps"`
	SwitchOn bool `yaml:"switch_on"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type RegisterConf struct {
	SwitchOn bool `yaml:"switch_on"`
}

type LogConf struct {
	Level string `yaml:"level"`
	Dir   string `yaml:"path"`
}

func initDir(serviceName string) (err error) {

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
	koalaConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIdx]), "..")
	koalaConf.ConfigDir = path.Join(koalaConf.RootDir, "./conf/", util.GetEnv())
	koalaConf.ConfigFile = path.Join(koalaConf.ConfigDir, fmt.Sprintf("%s.yaml", serviceName))
	return
}

func InitConfig(serviceName string) (err error) {

	err = initDir(serviceName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(koalaConf.ConfigFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &koalaConf)
	if err != nil {
		return
	}

	fmt.Printf("init koala conf succ, conf:%#v\n", koalaConf)
	return
}

func GetConfigDir() string {
	return koalaConf.ConfigDir
}

func GetRootDir() string {
	return koalaConf.RootDir
}

func GetServerPort() int {
	return koalaConf.Port
}

func GetConf() *KoalaConf {
	return koalaConf
}
