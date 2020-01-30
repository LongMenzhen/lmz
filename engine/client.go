package engine

import (
	"fmt"
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
	Group *Group          // 客户端所属房间
	Conn  *websocket.Conn // 客户端连接
	Send  chan []byte     // 待发送给客户端的内容
	Done  chan bool       // 是否连接已结束
}

// ReadMessage 读消息
func (c *Client) ReadMessage() {
	defer func() {
		c.Group.RemoveClient(c)
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
			log.Println("解析上行数据Json格式失败" + err.Error())
			continue
		}
		ctx := NewContext(*message)

		go func(ctx *Context) {
			for _, action := range Actions {
				if action.Event == message.Event {
					fmt.Println("0000000000")
					action.F(ctx)
					fmt.Println("333333")
					break
				}
			}
		}(ctx)

		fmt.Println("1111111111")
		result := <-ctx.Response
		fmt.Println(string(result))
		// data = bytes.TrimSpace(bytes.Replace(data, newline, space, -1))
		c.Group.Broadcast <- result
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

	gid := r.URL.Query().Get("group_id")
	if "" == gid {
		panic("请指定连接组")
		// http.Error(w, "请指定连接房间", 403)
		// return
	}

	groupID, _ := strconv.Atoi(gid)
	group := hub.GroupByID(int32(groupID))

	// 连接房间不存在
	if nil == group {
		panic("连接组不存在")
		// http.Error(w, "连接房间不存在", 422)
		// return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		http.Error(w, "升级websocket协议失败", 403)
		return
	}
	client := &Client{Group: group, Conn: conn, Send: make(chan []byte, 512), Done: make(chan bool)}
	group.AddClient(client)

	go client.WriteMessage()
	go client.ReadMessage()
}
