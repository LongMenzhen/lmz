package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config 配置
var Config Setting

// Setting 配置结构体
type Setting struct {
	Redis RedisSetting `yaml:"redis"`
}

// RedisSetting Redis配置
type RedisSetting struct {
	Host string `yaml:"hose"`
	Port int    `yaml:"port"`
}

func init() {
	parse()
}

func parse() {
	openFile, err := os.OpenFile("./config/config.yml", os.O_RDWR, 0644)
	source, err := ioutil.ReadAll(openFile)
	if nil != err {
		panic("配置文件读取失败: " + err.Error())
	}

	if err := yaml.Unmarshal(source, &Config); nil != err {
		panic("Yaml配置文件反序列化失败: " + err.Error())
	}
}
