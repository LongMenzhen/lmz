package engine

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 6 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

// Upgrader 升级
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 服务端注册客户端结构体
// 这里设计是1个客户端属于一个房间
// 后续可以升级为 N 对 M  的设计
type Client struct {
	Groups map[*Group]bool // 客户端所属房间
	Conn   *websocket.Conn // 客户端连接
	Send   chan []byte     // 待发送给客户端的内容
	Done   chan bool       // 是否连接已结束
}

// ReadMessage 读消息
func (c *Client) ReadMessage() {
	defer func() {
		for group := range c.Groups {
			group.RemoveClient(c)
		}
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if nil != err {
			log.Println("读取客户端信息失败", err.Error())
			return
		}

		message, err := NewMessage(data)
		if nil != err {
			log.Println("客户端上行原始数据JSON格式反序列化失败: " + err.Error())
			continue
		}

		done := make(chan bool, 1)
		ctx := NewContext(*message, c)

		// 不在消息组内，那么不能给该组发送消息
		if nil == ctx.Group {
			continue
		}

		go func(ctx Context) {
			for _, action := range Actions {
				if action.Event == message.Event {
					action.F(ctx)
					done <- true
					break
				}
			}
		}(ctx)

		// 处理F 没有返回数据的情况
		select {
		case <-done:
			log.Println("<-done")
			close(ctx.Response)
		case result := <-ctx.Response:
			log.Println("<-result")
			ctx.Group.Broadcast <- result
			close(done)
		}
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
