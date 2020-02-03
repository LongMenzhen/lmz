package route

import (
	"net/http"

	"github.com/cyrnicolase/lmz/app"
	"github.com/cyrnicolase/lmz/app/ws"
	"github.com/cyrnicolase/lmz/engine"
)

func init() {
	// 注册ws 接收消息事件
	engine.Registe("somebody", ws.SomebodyAction)
	engine.Registe("welcome", ws.WelcomeAction)
	engine.Registe("ping", ws.PingAction)
	engine.Registe("join-group", ws.JoinGroupAction)
}

// Route 注册路由
func Route() {
	http.HandleFunc("/ws", app.ServeWs)
	http.HandleFunc("/group", app.BuildGroup)
}
