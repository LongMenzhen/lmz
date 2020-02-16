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
	Socks      map[SockID]*Sock
	Register   chan *Sock
	Unregister chan *Sock
	Broadcast  chan []byte
}

// AddSock 将客户端连接加入到仓库
func (h *Hub) AddSock(sock *Sock) {
	h.Socks[sock.ID] = sock
}

// RemoveSock 从仓库中移除客户端
func (h *Hub) RemoveSock(sock *Sock) {
	if _, ok := h.Socks[sock.ID]; ok {
		sock.Done <- true
		delete(h.Socks, sock.ID)
	}
}

// Run 运行仓库监听
func (h *Hub) Run() {
	// 同步观察，后续是否有新的组进来，如果有新组创建或者结束
	go func() {
		for {
			select {
			case sock := <-h.Register:
				h.AddSock(sock)
			case sock := <-h.Unregister:
				h.RemoveSock(sock)
			case message := <-h.Broadcast:
				for _, sock := range h.Socks {
					sock.Send <- message
				}
			}
		}
	}()
}

// AttachHub 新建服务端仓库,如果仓库存在，那么直接返回该仓库
func AttachHub() *Hub {
	onceFunc := func() {
		hub = &Hub{
			Socks:      make(map[SockID]*Sock),
			Register:   make(chan *Sock),
			Unregister: make(chan *Sock),
			Broadcast:  make(chan []byte),
		}
	}
	once.Do(onceFunc)

	return hub
}
