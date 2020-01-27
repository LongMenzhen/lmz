package entity

var roomID int32 = 0

// Room 聊天房间
type Room struct {
	ID        int32            // 房间id
	Clients   map[*Client]bool // 房间拥有客户端
	Broadcast chan []byte      // 广播信息
	Done      chan bool        // 房间结束,收到这个结束信号，那么就关闭房间运行
}

// NewRoom 创建房间
func NewRoom() *Room {
	roomID++
	return &Room{
		ID:        roomID,
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte, 256),
		Done:      make(chan bool),
	}
}

// AddClient 将客户端加入到房间
func (r *Room) AddClient(client *Client) {
	if _, ok := r.Clients[client]; !ok {
		r.Clients[client] = true
	}
}

// RemoveClient 移除客户端
func (r *Room) RemoveClient(client *Client) {
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
	}
}

// KillClients 断开房间与客户端的关联关系
func (r *Room) KillClients() {
	for client := range r.Clients {
		r.RemoveClient(client)
	}
}

// Run 运行房间监听
func (r *Room) Run() {
	for {
		select {
		case <-r.Done: // 房间已完结
			r.KillClients()
			return
		case message := <-r.Broadcast:
			for client := range r.Clients {
				client.Send <- message
			}
		}
	}
}
