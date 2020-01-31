package engine

import "encoding/json"

// Context 操作上下文
type Context struct {
	Request  Message
	Response chan []byte
}

// NewContext 新生成Context结构体
func NewContext(msg Message) Context {
	return Context{
		Request:  msg,
		Response: make(chan []byte),
	}
}

// String 字符串输出
func (ctx *Context) String(str string) {
	ctx.Response <- []byte(str)
}

// Bytes 字符流输出
func (ctx *Context) Bytes(bytes []byte) {
	ctx.Response <- bytes
}

// JSON json格式输出
func (ctx *Context) JSON(j map[string]interface{}) {
	bytes, _ := json.Marshal(j)
	ctx.Bytes(bytes)
}
