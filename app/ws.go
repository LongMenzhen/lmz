package app

import (
	"net/http"

	"github.com/cyrnicolase/lmz/engine"
)

// ServeWs 提供Websocket 服务
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := engine.Upgrader.Upgrade(w, r, nil)
	if nil != err {
		http.Error(w, "升级websocket协议失败", 403)
		return
	}
	client := engine.NewClient(conn)

	go client.WriteMessage()
	go client.ReadMessage()
}
