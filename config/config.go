package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 配置
var Config Setting

// Setting 配置结构体
type Setting struct {
	HTTP  httpSetting  `yaml:"http"`
	Redis redisSetting `yaml:"redis"`
}

type httpSetting struct {
	Addr string `yaml:"addr"`
}

// redisSetting Redis配置
type redisSetting struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func init() {
	parse()
}

func parse() {
	openFile, err := os.OpenFile("./config/config.yml", os.O_RDONLY, 0644)
	source, err := ioutil.ReadAll(openFile)
	if nil != err {
		panic("配置文件读取失败: " + err.Error())
	}

	if err := yaml.Unmarshal(source, &Config); nil != err {
		panic("Yaml配置文件反序列化失败: " + err.Error())
	}
}
