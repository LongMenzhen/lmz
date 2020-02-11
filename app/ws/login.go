package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/model"
	// "github.com/vmihailenco/msgpack/v4"
)

var (
	space       = []byte{' '}
	newline     = []byte{'\n'}
	htmlNewline = []byte("<br />")
	htmlSpace   = []byte("&nbsp;")
)

var users map[int32]UserMap = map[int32]UserMap{}

// UserMap 用户信息
type UserMap struct {
	Name   string
	Client *engine.Client
}

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

	user := UserMap{
		Name:   request.Name,
		Client: ctx.Client,
	}

	// 如果用户已经存在
	for _, u := range users {
		if user.Name == u.Name {
			return
		}
	}

	// TODO 用户名唯一性校验
	users[user.Client.ID] = user
	data := ctx.Format(map[string]interface{}{
		"name": user.Name,
	})

	hub := engine.AttachHub()
	hub.Broadcast <- data

	// b, err := msgpack.Marshal(&user)
	// if nil != err {
	// 	ctx.String("保存登陆用户信息失败")
	// 	return
	// }

	key := fmt.Sprintf("login:user:%d", user.Client.ID)
	log.Println("=======", key, user.Name)
	redis := model.Redis()
	redis.Set(key, user.Name, 60*time.Second)

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

	// content := bytes.Replace([]byte(request.Content), newline, htmlNewline, -1)
	// content = bytes.Replace(content, space, htmlSpace, -1)
	user := users[ctx.Client.ID]
	name := user.Name

	type Result struct {
		Name      string    `json:"name"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	result := Result{
		Name:      name,
		Content:   string(request.Content),
		CreatedAt: time.Now(),
	}

	ctx.Mix(result)
}
