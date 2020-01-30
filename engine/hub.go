package engine

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
// 各个客户端在选择连接房间的时候，携带 group_id
// 仓库只能够初始化一次，所以要设计为单例
type Hub struct {
	Groups     map[int32]*Group // 仓库房间集合
	Register   chan *Group      // 新建房间
	Unregister chan *Group      // 销毁房间
}

// AddGroup 将房间加入到仓库
func (h *Hub) AddGroup(group *Group) {
	go group.Run()
	h.Groups[group.ID] = group
}

// RemoveGroup 移除房间
func (h *Hub) RemoveGroup(group *Group) {
	if _, ok := h.Groups[group.ID]; ok {
		group.Done <- true         // 关闭房间服务
		delete(h.Groups, group.ID) // 从集合中删除该房间
	}
}

// GroupByID 根据房间id查询仓库房间
func (h *Hub) GroupByID(groupID int32) *Group {
	if group, ok := h.Groups[groupID]; ok {
		return group
	}

	return nil
}

// Run 运行仓库监听
func (h *Hub) Run() {
	// 同步观察，后续是否有新的房间进来，如果有新房间创建或者结束
	go func() {
		for {
			select {
			case group := <-h.Register: // http 创建一个新房间
				h.AddGroup(group)
			case group := <-h.Unregister: // 某个操作结束后，销毁房间
				h.RemoveGroup(group)
			}
		}
	}()
}

// AttachHub 新建服务端仓库,如果仓库存在，那么直接返回该仓库
func AttachHub() *Hub {
	onceFunc := func() {
		hub = &Hub{
			Groups:     make(map[int32]*Group),
			Register:   make(chan *Group),
			Unregister: make(chan *Group),
		}
	}
	once.Do(onceFunc)

	return hub
}
