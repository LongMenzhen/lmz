package model

import (
	"github.com/cyrnicolase/lmz/engine"
)

// SockMap 客户端连接维护
// key 为 user_id
// value 为连接客户端
var SockMap = map[UserID]*engine.Sock{}

// AddUserSock 将登陆用户与客户端socket关联
func AddUserSock(userID UserID, client *engine.Sock) {
	SockMap[userID] = client
}
