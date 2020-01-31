package app

import (
	"encoding/json"

	"github.com/cyrnicolase/lmz/engine"
)

// SomebodyAction 某人说
func SomebodyAction(ctx engine.Context) {
	type Request struct {
		Name string `json:"name"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		ctx.String("格式化上行数据失败")
		return
	}

	ctx.String("hello: " + request.Name)
}

// WelcomeAction 欢迎
func WelcomeAction(ctx engine.Context) {
	type Request struct {
		Name string `json:"name"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		ctx.String("格式化上行数据失败")
		return
	}

	ctx.Bytes([]byte("Welcome: " + request.Name))
}
