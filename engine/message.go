package engine

import (
	"encoding/json"
	"log"
)

// Message 消息体格式
// 规定客户端与服务端交互的消息体格式
// {
//     event:(string),
// 	   group_id:(int),
//     body: {
//
//     }
// }
type Message struct {
	Event   string          `json:"event"`
	GroupID int32           `json:"group_id"`
	Body    json.RawMessage `json:"body"`
}

// NewMessage 根据给定字节流反序列化为Message结构体
func NewMessage(src []byte) (*Message, error) {
	var message Message
	log.Println("input message: " + string(src))
	err := json.Unmarshal(src, &message)
	if nil != err {
		return nil, err
	}

	return &message, nil
}
