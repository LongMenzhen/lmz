package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/cyrnicolase/lmz/model"
	"github.com/go-redis/redis"
)

// LuaDir Lua脚本目录
const LuaDir = "./lua/"

// 将Lua脚本载入到Redis

func main() {
	rds := model.Redis()
	rds.ScriptFlush()
	names := iter()
	sha1s := map[string]string{}
	for _, name := range names {
		bytea, err := ioutil.ReadFile(name)
		if nil != err {
			panic(err)
		}

		sc := redis.NewScript(string(bytea))
		sha1 := sc.Load(rds).String()

		sha1s[name] = sha1
	}

	fmt.Println(sha1s)
}

// iter 遍历目录
func iter() []string {
	files, err := ioutil.ReadDir(LuaDir)
	if nil != err {
		panic("读取lua脚本失败")
	}

	names := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if ".lua" == path.Ext(filename) {
			name := LuaDir + filename
			names = append(names, name)
		}
	}

	return names
}
