package ws

import (
	"encoding/json"
	"fmt"

	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/model"
	// "github.com/vmihailenco/msgpack/v4"
)

var users map[*engine.Client]model.User = map[*engine.Client]model.User{}

// LoginAction 登陆
// 第一次进入房间，要求必须有全局唯一的用户信息
// 这里使用用户名做全局唯一控制
func LoginAction(ctx engine.Context) {
	type Request struct {
		Name string `json:"name"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		ctx.String("登陆信息缺少登陆用户名")
		return
	}

	user := model.User{
		Name:     request.Name,
		ClientID: ctx.Request.GroupID,
	}

	users[ctx.Client] = user

	// b, err := msgpack.Marshal(&user)
	// if nil != err {
	// 	ctx.String("保存登陆用户信息失败")
	// 	return
	// }

	// key := fmt.Sprintf("login:user:%d", user.ClientID)
	// redis := model.Redis()
	// redis.Set(key, b)

}

// SayAction 说点什么
func SayAction(ctx engine.Context) {
	type Request struct {
		Content string `json:"content"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		ctx.String("发送消息格式错误" + err.Error())
		return
	}

	user := users[ctx.Client]
	name := user.Name

	ctx.String(fmt.Sprintf("%s说:%s", name, request.Content))
}
