package entity

var roomID int32 = 0

// Room 聊天房间
type Room struct {
	ID        int32            `json:"id"`        // 房间id
	Clients   map[*Client]bool `json:"clients"`   // 房间拥有客户端
	Broadcast chan []byte      `json:"broadcast"` // 广播信息
}

// NewRoom 创建房间
func NewRoom() *Room {
	roomID++
	return &Room{
		ID:        roomID,
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte, 512),
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

// Run 运行房间监听
func (r *Room) Run() {
	for {
		select {
		case message := <-r.Broadcast:
			for client := range r.Clients {
				client.Send <- message
			}
		}
	}
}
