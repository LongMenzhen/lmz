package engine

import (
	"log"
	"net/http"
)

// HandlerFunc 执行函数
type HandlerFunc func(ctx Context)

// Action 操作
type Action struct {
	Event string
	F     HandlerFunc
}

// Actions 注册执行方法集合
var Actions []Action

// Registe 注册
func Registe(event string, f HandlerFunc) {
	Actions = append(Actions, Action{
		Event: event,
		F:     f,
	})
}

// RunServe 运行服务
func RunServe(addr string) {
	if err := http.ListenAndServe(addr, nil); nil != err {
		log.Fatal("监听失败" + err.Error())
	}
}
