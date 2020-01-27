package entity

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 6 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client 服务端注册客户端结构体
// 这里设计是1个客户端属于一个房间
// 后续可以升级为 N 对 M  的设计
type Client struct {
	Room *Room           // 客户端所属房间
	Conn *websocket.Conn // 客户端连接
	Send chan []byte     // 待发送给客户端的内容
	Done chan bool       // 是否连接已结束
}

// ReadMessage 读消息
func (c *Client) ReadMessage() {
	defer func() {
		c.Room.RemoveClient(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if nil != err {
			log.Println("读取客户端信息失败", err.Error())
			return
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Room.Broadcast <- message
	}
}

// WriteMessage 发消息
func (c *Client) WriteMessage() {
	tick := time.NewTicker(time.Second)
	defer func() {
		tick.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case <-c.Done:
			return
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if nil != err {
				return
			}
			w.Write(message)
			if err := w.Close(); nil != err {
				return
			}

		case <-tick.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); nil != err {
				return
			}
		}
	}
}

// ServeWs 提供websocket服务
func ServeWs(w http.ResponseWriter, r *http.Request) {
	hub := AttachHub()
	rid := r.URL.Query().Get("room_id")
	if "" == rid {
		http.Error(w, "请指定连接房间", 403)
		return
	}

	roomID, _ := strconv.Atoi(rid)
	room := hub.RoomByID(int32(roomID))

	// 连接房间不存在
	if nil == room {
		http.Error(w, "连接房间不存在", 422)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		http.Error(w, "升级websocket协议失败", 403)
		return
	}

	client := &Client{Room: room, Conn: conn, Send: make(chan []byte, 512), Done: make(chan bool)}
	room.AddClient(client)

	go client.WriteMessage()
	go client.ReadMessage()
}
