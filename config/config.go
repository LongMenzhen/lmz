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
	Luas  luasSetting  `yaml:"luas"`
}

type httpSetting struct {
	Addr string `yaml:"addr"`
}

// redisSetting Redis配置
type redisSetting struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// luasSetting lua操作Redis的脚本
type luasSetting struct {
	MultGetGroupNames string `yaml:"mult_get_group_names"`
	MultGroups        string `yaml:"mult_groups"`
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
