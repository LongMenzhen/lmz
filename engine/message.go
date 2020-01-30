package engine

import "encoding/json"

// Message 消息体格式
// 规定客户端与服务端交互的消息体格式
// {
//     event:(string),
//     body: {
//
//     }
// }
type Message struct {
	Event string
	Body  json.RawMessage
}

// NewMessage 根据给定字节流反序列化为Message结构体
func NewMessage(src []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(src, &message)
	if nil != err {
		return nil, err
	}

	return &message, nil
}
