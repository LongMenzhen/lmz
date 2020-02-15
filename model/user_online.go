package model

import (
	"fmt"

	"github.com/cyrnicolase/lmz/config"
)

// ClientUser socket客户端与用户关联关系
type ClientUser struct {
	ClientID int32
	UserID   int32
}

// TableName 返回Rediskey
func (ClientUser) TableName(clientID int32) string {
	return fmt.Sprintf("user:client:%d", int(clientID))
}

// NewClientUser 创建客户端、顾客之间关系
func NewClientUser(clientID, userID int32) *ClientUser {
	return &ClientUser{
		ClientID: clientID,
		UserID:   userID,
	}
}

// ClientIDs 在线客户端连接
type ClientIDs []int32

// TableName 返回Redis key
func (ClientIDs) TableName() string {
	return "clientids"
}

// CreateClientUser 创建客户端用户关联关系
func CreateClientUser(clientUser ClientUser) error {
	pipe := rds.Pipeline()
	defer pipe.Close()
	err := pipe.Set(clientUser.TableName(clientUser.ClientID), int32(clientUser.UserID), 0).Err()
	if nil != err {
		return err
	}
	pipe.SAdd(ClientIDs{}.TableName(), clientUser.ClientID)
	_, err = pipe.Exec()
	if nil != err {
		return err
	}

	return nil
}

// MultGetNames 返回当前登录的所有用户名
func MultGetNames() []string {
	ret, _ := rds.EvalSha(config.Config.Luas.MultGetNames, []string{ClientIDs{}.TableName()}, nil).Result()
	names, ins := []string{}, ret.([]interface{})

	for _, inter := range ins {
		names = append(names, inter.(string))
	}

	return names
}
