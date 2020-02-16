package model

import (
	"github.com/cyrnicolase/lmz/engine"
	"github.com/sirupsen/logrus"
)

// SockMap 客户端连接维护
// key 为 user_id
// value 为连接客户端
var SockMap = map[UserID]*engine.Sock{}

// AddUserSock 将登陆用户与客户端socket关联
func AddUserSock(userID UserID, client *engine.Sock) {
	logrus.WithFields(map[string]interface{}{
		"user_id":   userID,
		"client_id": client.ID,
	}).Info("绑定用户到客户端")

	SockMap[userID] = client
}
