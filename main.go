package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/cyrnicolase/lmz/engine"
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
		BuildRoom(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); nil != err {
		log.Fatal("监听失败" + err.Error())
	}
}

// BuildRoom 创建房间，并注册房间的Hub
func BuildRoom(w http.ResponseWriter, r *http.Request) {
	hub := engine.AttachHub()
	if "GET" == r.Method {
		rooms := hub.Groups
		i, roomIDs := 0, make([]int32, len(rooms))
		for roomID := range rooms {
			roomIDs[i] = roomID
			i++
		}

		data, _ := json.Marshal(roomIDs)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
		return
	}

	newRoom := engine.NewGroup()
	hub.Register <- newRoom

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("房间创建完成: " + strconv.Itoa(int(newRoom.ID))))
}
