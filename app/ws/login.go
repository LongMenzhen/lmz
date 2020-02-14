package ws

import (
	"encoding/json"

	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// LoginAction 登陆
func LoginAction(ctx engine.Context) {
	type Request struct {
		UserID   int32
		Password string
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		logrus.Error("ws登陆用户登陆信息json解析失败" + err.Error())
		ctx.Error("ws登陆用户信息json解释错误")
		return
	}

	user := &model.User{ID: model.UserID(request.UserID)}
	if err := user.MakeUser(); nil != err {
		logrus.Error("登陆用户不存在" + err.Error())
		ctx.Error("登陆用户不存在")
		return
	}

	if user.Password != request.Password {
		logrus.Error("登陆账号密码错误")
		ctx.Error("登陆账号密码错误")
		return
	}

	// 关联clientID 与 userID
	clientID := ctx.Client.ID
	clientUser := model.NewClientUser(clientID, user.ID)
	if err := model.CreateClientUser(*clientUser); nil != err {
		logrus.Error("关联ws客户端id与顾客id失败" + err.Error())
		ctx.Error("登陆失败")
		return
	}

	ctx.String("登陆成功")
}

// SayAction 说点什么
func SayAction(ctx engine.Context) {

}
