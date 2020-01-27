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
	Rooms      map[int32]*Room // 仓库房间集合
	Register   chan *Room      // 新建房间
	Unregister chan *Room      // 销毁房间
}

// AddRoom 将房间加入到仓库
func (h *Hub) AddRoom(room *Room) {
	go room.Run()
	h.Rooms[room.ID] = room
}

// RemoveRoom 移除房间
func (h *Hub) RemoveRoom(room *Room) {
	if _, ok := h.Rooms[room.ID]; ok {
		room.Done <- true        // 关闭房间服务
		delete(h.Rooms, room.ID) // 从集合中删除该房间
	}
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
	// 同步观察，后续是否有新的房间进来，如果有新房间创建或者结束
	go func() {
		for {
			select {
			case room := <-h.Register: // http 创建一个新房间
				h.AddRoom(room)
			case room := <-h.Unregister: // 某个操作结束后，销毁房间
				h.RemoveRoom(room)
			}
		}
	}()
}

// AttachHub 新建服务端仓库,如果仓库存在，那么直接返回该仓库
func AttachHub() *Hub {
	onceFunc := func() {
		hub = &Hub{
			Rooms:      make(map[int32]*Room),
			Register:   make(chan *Room),
			Unregister: make(chan *Room),
		}
	}
	once.Do(onceFunc)

	return hub
}
