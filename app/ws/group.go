package ws

import (
	"strconv"

	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// CreateGroup 创建组
func CreateGroup(ctx engine.Context) {
	user, _ := model.MakeUserByClientID(ctx.Client.ID)
	group := model.NewGroup(user)
	if err := model.CreateGroup(*group); nil != err {
		logrus.Error("创建组失败")
		ctx.Error("创建消息组失败")
		return
	}

	ctx.String("创建组成功:" + strconv.Itoa(int(group.ID)))
}
