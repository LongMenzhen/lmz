package ws

import "github.com/cyrnicolase/lmz/engine"

// JoinGroupAction 加入组
// 在对组发消息之前，都应该先将该客户端加入到组中
func JoinGroupAction(ctx engine.Context) {
	group := ctx.Group
	client := ctx.Client

	group.AddClient(client)
}
