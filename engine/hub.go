package engine

import (
	"sync"
)

var (
	hub  *Hub
	once sync.Once
)

// Hub 服务端组仓库，管理所有的组
// 先创建好组后，将组加入到Hub 中
// 然后就将创建好的组ID 通知到各个客户端
// 各个客户端在选择连接组的时候，携带 group_id
// 仓库只能够初始化一次，所以要设计为单例
type Hub struct {
	Clients    map[int32]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

// AddClient 将客户端连接加入到仓库
func (h *Hub) AddClient(client *Client) {
	h.Clients[client.ID] = client
}

// RemoveClient 从仓库中移除客户端
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.Clients[client.ID]; ok {
		client.Done <- true
		delete(h.Clients, client.ID)
	}
}

// Run 运行仓库监听
func (h *Hub) Run() {
	// 同步观察，后续是否有新的组进来，如果有新组创建或者结束
	go func() {
		for {
			select {
			case client := <-h.Register:
				h.AddClient(client)
			case client := <-h.Unregister:
				h.RemoveClient(client)
			case message := <-h.Broadcast:
				for _, client := range h.Clients {
					client.Send <- message
				}
			}
		}
	}()
}

// AttachHub 新建服务端仓库,如果仓库存在，那么直接返回该仓库
func AttachHub() *Hub {
	onceFunc := func() {
		hub = &Hub{
			Clients:    make(map[int32]*Client),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
			Broadcast:  make(chan []byte),
		}
	}
	once.Do(onceFunc)

	return hub
}
