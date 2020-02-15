package model

import (
	"fmt"

	"github.com/cyrnicolase/lmz/config"
)

// UserClient socket客户端与用户关联关系
type UserClient struct {
	ClientID int32
	UserID   int32
}

// TableName 返回Rediskey
// user:client:{user_id}
func (UserClient) TableName(userID int32) string {
	return fmt.Sprintf("user:client:%d", int(userID))
}

// NewUserClient 创建客户端、顾客之间关系
func NewUserClient(clientID, userID int32) *UserClient {
	return &UserClient{
		ClientID: clientID,
		UserID:   userID,
	}
}

// UserOnlines 在线客户端连接
type UserOnlines []int32

// TableName 返回Redis key
func (UserOnlines) TableName() string {
	return "user:online:userids"
}

// CreateUserClient 创建客户端用户关联关系
func CreateUserClient(uc UserClient) error {
	pipe := rds.Pipeline()
	defer pipe.Close()
	err := pipe.Set(uc.TableName(uc.ClientID), int32(uc.UserID), 0).Err()
	if nil != err {
		return err
	}
	pipe.SAdd(UserOnlines{}.TableName(), uc.ClientID)
	_, err = pipe.Exec()
	if nil != err {
		return err
	}

	return nil
}

// MultGetNames 返回当前登录的所有用户名
func MultGetNames() []string {
	ret, _ := rds.EvalSha(config.Config.Luas.MultGetNames, []string{UserOnlines{}.TableName()}, nil).Result()
	names, ins := []string{}, ret.([]interface{})

	for _, inter := range ins {
		names = append(names, inter.(string))
	}

	return names
}
