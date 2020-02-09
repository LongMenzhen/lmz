package app

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// BuildGroup 创建组
func BuildGroup(w http.ResponseWriter, r *http.Request) {
	hub := engine.AttachHub()
	if "GET" == r.Method {
		groups := hub.Groups
		i, groupIDs := 0, make([]int32, len(groups))
		for groupID := range groups {
			groupIDs[i] = groupID
			i++
		}

		data, _ := json.Marshal(groupIDs)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	newGroup := engine.NewGroup()
	hub.Register <- newGroup
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Websocket消息组创建完成: group_id=" + strconv.Itoa(int(newGroup.ID))))
}
