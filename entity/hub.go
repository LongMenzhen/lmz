package entity

import (
	"sync"
)

var (
	hub  *Hub
	once sync.Once
)

// Hub 服务端房间仓库，管理所有的房间
// 先创建好房间后，将房间加入到Hub 中
// 然后就将创建好的房间ID 通知到各个客户端
// 各个客户端在选择连接房间的时候，携带 room_id
// 仓库只能够初始化一次，所以要设计为单例
type Hub struct {
	Rooms map[int32]*Room
}

// AddRoom 将房间加入到仓库
func (h *Hub) AddRoom(room *Room) {
	h.Rooms[room.ID] = room
}

// RoomByID 根据房间id查询仓库房间
func (h *Hub) RoomByID(roomID int32) *Room {
	if room, ok := h.Rooms[roomID]; ok {
		return room
	}

	return nil
}

// Run 运行仓库监听
func (h *Hub) Run() {
	for _, room := range h.Rooms {
		go room.Run() // 使用协程一直运行
	}
}

// NewHub 新建服务端仓库
func NewHub() *Hub {
	onceFunc := func() {
		hub = &Hub{
			Rooms: make(map[int32]*Room, 0),
		}
	}
	once.Do(onceFunc)

	return hub
}
