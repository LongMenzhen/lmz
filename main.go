package main

import (
	"log"
	"net/http"

	"github.com/cyrnicolase/lmz/entity"
)

func main() {

	room := entity.NewRoom()
	room1 := entity.NewRoom()
	hub := entity.NewHub()
	hub.AddRoom(room)
	hub.AddRoom(room1)
	hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		entity.ServeWs(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); nil != err {
		log.Fatal("监听失败" + err.Error())
	}
}
