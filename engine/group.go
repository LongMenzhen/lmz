package engine

import (
	"log"
	"sync"
)

// GroupID 组id；考虑到并发的情况，这里使用锁
type GroupID struct {
	ID int32
	m  *sync.RWMutex
}

var groupID GroupID

// IncrGroupID 自增组id
func IncrGroupID() int32 {
	groupID.m.Lock()
	defer groupID.m.Unlock()

	groupID.ID++

	return groupID.ID
}

// Group 聊天房间
type Group struct {
	ID        int32            // 房间id
	Clients   map[*Client]bool // 房间拥有客户端
	Broadcast chan []byte      // 广播信息
	Done      chan bool        // 房间结束,收到这个结束信号，那么就关闭房间运行
}

// NewGroup 创建房间
func NewGroup() *Group {
	id := IncrClientID()
	return &Group{
		ID:        id,
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte, 256),
		Done:      make(chan bool),
	}
}

// AddClient 将客户端加入到房间
func (r *Group) AddClient(client *Client) {
	if _, ok := r.Clients[client]; !ok {
		r.Clients[client] = true
	}
}

// HasClient 判断组内是否存指定客户端
func (r *Group) HasClient(client *Client) bool {
	if _, ok := r.Clients[client]; !ok {
		return false
	}

	return true
}

// RemoveClient 移除客户端
func (r *Group) RemoveClient(client *Client) {
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
	}
}

// KillClients 断开房间与客户端的关联关系
func (r *Group) KillClients() {
	for client := range r.Clients {
		r.RemoveClient(client)
	}
}

// Run 运行房间监听
func (r *Group) Run() {
	for {
		select {
		case <-r.Done: // 房间已完结
			r.KillClients()
			return
		case message := <-r.Broadcast:
			log.Println("output message:" + string(message))
			for client := range r.Clients {
				client.Send <- message
			}
		}
	}
}
