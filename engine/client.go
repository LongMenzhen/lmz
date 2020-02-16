package engine

import (
	"log"
	"net/http"
	"sync"
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

// ClientID 客户端id
type ClientID int32

// RwClientID 这里需要考虑并发问题，所以要进行锁操作
type RwClientID struct {
	ID ClientID
	m  *sync.RWMutex
}

var newClientID = RwClientID{
	ID: 0,
	m:  &sync.RWMutex{},
}

// Client 服务端注册客户端结构体
// 这里设计是1个客户端属于一个房间
// 后续可以升级为 N 对 M  的设计
type Client struct {
	ID   ClientID
	Conn *websocket.Conn // 客户端连接
	Send chan []byte     // 待发送给客户端的内容
	Done chan bool       // 是否连接已结束
}

// ReadMessage 读消息
func (c *Client) ReadMessage() {
	defer func() {
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
		noEvent := make(chan bool, 1)
		ctx := NewContext(*message, c)

		// 不在消息组内，那么不能给该组发送消息
		// if nil == ctx.Group {
		// 	continue
		// }

		go func(ctx Context) {
			for _, action := range Actions {
				if action.Event == message.Event {
					action.F(ctx)
					done <- true
					return
				}
			}
		}(ctx)

		// 处理F 没有返回数据的情况
		select {
		case <-noEvent:
			log.Println("<-noEvent")
			close(ctx.Response)
			close(done)
		case <-done:
			log.Println("<-done")
			close(ctx.Response)
			close(noEvent)
		case result := <-ctx.Response:
			log.Println("<-result")
			c.Send <- result
			close(done)
			close(noEvent)
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

// NewClient 创建新的客户端连接
func NewClient(conn *websocket.Conn) *Client {
	id := IncrRwClientID()
	client := &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte),
		Done: make(chan bool),
	}

	// 如果没有登陆，就不能注册进入hub中
	// hub := AttachHub()
	// hub.RegisterClient <- client // 注册客户端

	return client
}

// IncrRwClientID 自增客户端id
func IncrRwClientID() ClientID {
	newClientID.m.Lock()
	defer newClientID.m.Unlock()
	newClientID.ID++

	return newClientID.ID
}
