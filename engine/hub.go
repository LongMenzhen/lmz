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
	Groups     map[int32]*Group // 仓库组集合
	Register   chan *Group      // 新建组
	Unregister chan *Group      // 销毁组

	Clients          map[int32]*Client
	RegisterClient   chan *Client
	UnregisterClient chan *Client
	Broadcast        chan []byte
}

// AddGroup 将组加入到仓库
func (h *Hub) AddGroup(group *Group) {
	go group.Run()
	h.Groups[group.ID] = group
}

// RemoveGroup 移除组
func (h *Hub) RemoveGroup(group *Group) {
	if _, ok := h.Groups[group.ID]; ok {
		group.Done <- true         // 关闭组服务
		delete(h.Groups, group.ID) // 从集合中删除该组
	}
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

// GroupByID 根据组id查询仓库组
func (h *Hub) GroupByID(groupID int32) *Group {
	if 0 == groupID {
		return nil
	}

	if group, ok := h.Groups[groupID]; ok {
		return group
	}

	return nil
}

// Run 运行仓库监听
func (h *Hub) Run() {
	// 同步观察，后续是否有新的组进来，如果有新组创建或者结束
	go func() {
		for {
			select {
			case group := <-h.Register: // http 创建一个新组
				h.AddGroup(group)
			case group := <-h.Unregister: // 某个操作结束后，销毁组
				h.RemoveGroup(group)
			case client := <-h.RegisterClient:
				h.AddClient(client)
			case client := <-h.UnregisterClient:
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
			Groups:           make(map[int32]*Group),
			Register:         make(chan *Group),
			Unregister:       make(chan *Group),
			Clients:          make(map[int32]*Client),
			RegisterClient:   make(chan *Client),
			UnregisterClient: make(chan *Client),
			Broadcast:        make(chan []byte),
		}
	}
	once.Do(onceFunc)

	return hub
}
