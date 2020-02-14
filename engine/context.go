package engine

import (
	"encoding/json"
)

// Response 下行数据
type Response struct {
	Event string      `json:"event"`
	Body  interface{} `json:"body"`
}

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
	ctx.Response <- ctx.Format(str)
}

// Bytes 字符流输出
func (ctx *Context) Bytes(bytea []byte) {
	ctx.Response <- ctx.Format(string(bytea))
}

// JSON json格式输出
func (ctx *Context) JSON(j map[string]interface{}) {
	// bytea, _ := json.Marshal(j)
	// ctx.Response <- ctx.Format(bytea)

	ctx.Response <- ctx.Format(j)
}

// Object 返回对象
func (ctx *Context) Object(o struct{}) {
	ctx.Response <- ctx.Format(o)
}

// Mix 混合模式
func (ctx *Context) Mix(m interface{}) {
	ctx.Response <- ctx.Format(m)
}

// Error 返回错误信息
func (ctx *Context) Error(msg string) {
	ctx.Request.Event = "error"
	ctx.Response <- ctx.Format(msg)
}

// Format 下行数据格式化
func (ctx Context) Format(i interface{}) []byte {
	resp := Response{
		Event: ctx.Request.Event,
		Body:  i,
	}

	data, _ := json.Marshal(resp)

	return data
}
