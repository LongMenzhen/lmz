package route

import (
	"net/http"

	"github.com/cyrnicolase/lmz/app"
	"github.com/cyrnicolase/lmz/app/ws"
	"github.com/cyrnicolase/lmz/engine"
)

func init() {
	// 注册ws 接收消息事件
	engine.Bind("somebody", ws.SomebodyAction)
	engine.Bind("welcome", ws.WelcomeAction)
	engine.Bind("ping", ws.PingAction)
	engine.Bind("login", ws.LoginAction)
	engine.Bind("join-group", ws.JoinGroupAction)
	engine.Bind("say", ws.SayAction)
}

// Route 注册路由
func Route() {
	http.HandleFunc("/ws", app.ServeWs)
	http.HandleFunc("/group", app.BuildGroup)
}
