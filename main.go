package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/cyrnicolase/lmz/engine"
	_ "github.com/cyrnicolase/lmz/route"
)

func main() {
	hub := engine.AttachHub()
	hub.Run()
	group := engine.NewGroup()
	group1 := engine.NewGroup()
	hub.AddGroup(group)
	hub.AddGroup(group1)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		engine.ServeWs(w, r)
	})

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		BuildGroup(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); nil != err {
		log.Fatal("监听失败" + err.Error())
	}
}

// BuildGroup 创建房间，并注册房间的Hub
func BuildGroup(w http.ResponseWriter, r *http.Request) {
	hub := engine.AttachHub()
	if "GET" == r.Method {
		groups := hub.Groups
		i, groupIDs := 0, make([]int32, len(groups))
		for roomID := range groups {
			groupIDs[i] = roomID
			i++
		}

		data, _ := json.Marshal(groupIDs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
		return
	}

	newRoom := engine.NewGroup()
	hub.Register <- newRoom

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("websocket 组创建完成: " + strconv.Itoa(int(newRoom.ID))))
}
