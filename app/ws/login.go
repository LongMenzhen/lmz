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
		UserID   model.UserID `json:"user_id,string"` // 这里标注上行的数据是string类型，从string类型转换成int32
		Password string       `json:"password"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		logrus.Error("ws登陆用户登陆信息json解析失败" + err.Error())
		ctx.Error("ws登陆用户信息json解释错误")
		return
	}

	user := &model.User{ID: request.UserID}
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
	clientUser := model.NewUserClient(clientID, user.ID)
	if err := model.CreateUserClient(*clientUser); nil != err {
		logrus.Error("关联ws客户端id与顾客id失败" + err.Error())
		ctx.Error("登陆失败")
		return
	}

	hub := engine.AttachHub()
	hub.Register <- ctx.Client
	names := map[string]interface{}{
		"names": model.MultGetNames(),
		"name":  user.Username,
	}

	// 将登陆用户与socket连接关联
	model.AddUserClient(user.ID, ctx.Client)

	ctx.Mix(names)
}

// JoinGroup 加入组
func JoinGroup(ctx engine.Context) {
	type Request struct {
		GroupID model.GroupID `json:"group_id"`
	}

	var request Request
	if err := json.Unmarshal(ctx.Request.Body, &request); nil != err {
		ctx.Error("解析加入组消息失败")
		return
	}

	groupID := request.GroupID
	group := &model.Group{ID: groupID}
	if err := group.MakeGroup(); nil != err {
		logrus.Error("待加入组不存在" + err.Error())
		ctx.Error("待加入组未找到")
		return
	}

	user, _ := model.MakeUserByClientID(ctx.Client.ID)
	model.AddUserToGroup(user, group)

	logrus.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"group_id": group.ID,
	}).Info("添加用户到组成功")

	ctx.String("添加用户到组成功")
}

// SayAction 说点什么
func SayAction(ctx engine.Context) {

}
