package ws

import (
	"encoding/json"

	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// CreateGroup 创建组
func CreateGroup(ctx engine.Context) {
	type Request struct {
		Name string `json:"name"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		logrus.Error("ws解析上行创建组消息失败" + err.Error())
		ctx.Error("解析上行消息格式失败")
		return
	}

	user, _ := model.MakeUserBySockID(ctx.Sock.ID)
	group := model.NewGroup(user, request.Name)
	if err := model.CreateGroup(*group); nil != err {
		logrus.Error("创建组失败")
		ctx.Error("创建消息组失败")
		return
	}

	ctx.Mix(group)
}

// GetGroups 返回消息组
func GetGroups(ctx engine.Context) {
	groups := model.MultGroups()
	ctx.Mix(groups)
}
