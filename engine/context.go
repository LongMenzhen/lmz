package engine

import "encoding/json"

// Context 操作上下文
type Context struct {
	Group    *Group
	Client   *Client
	Request  Message
	Response chan []byte
}

// NewContext 新生成Context结构体
func NewContext(msg Message, client *Client) Context {
	hub := AttachHub()
	group := hub.GroupByID(msg.GroupID)

	// 临时使用，实际上应该是客户端先加入组了才能进行组消息发送；现在发送组消息的时候，默认都将客户端加入进去
	if nil != group {
		group.AddClient(client)
	}

	return Context{
		Group:    group,
		Client:   client,
		Request:  msg,
		Response: make(chan []byte),
	}
}

// String 字符串输出
func (ctx *Context) String(str string) {
	ctx.Response <- []byte(str)
}

// Bytes 字符流输出
func (ctx *Context) Bytes(bytea []byte) {
	ctx.Response <- bytea
}

// JSON json格式输出
func (ctx *Context) JSON(j map[string]interface{}) {
	bytes, _ := json.Marshal(j)
	ctx.Bytes(bytes)
}
