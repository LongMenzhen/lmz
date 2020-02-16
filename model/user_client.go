package model

import (
	"fmt"

	"github.com/cyrnicolase/lmz/config"
	"github.com/cyrnicolase/lmz/engine"
)

// UserClient socket客户端与用户关联关系
type UserClient struct {
	ClientID engine.ClientID
	UserID   UserID
}

// TableName 返回Rediskey
// user:client:{user_id}
func (UserClient) TableName(clientid engine.ClientID) string {
	return fmt.Sprintf("user:client:%d", clientid)
}

// NewUserClient 创建客户端、顾客之间关系
func NewUserClient(clientID engine.ClientID, userID UserID) *UserClient {
	return &UserClient{
		ClientID: clientID,
		UserID:   userID,
	}
}

// MakeUserByClientID 根据客户端id创建用户对象
func MakeUserByClientID(cid engine.ClientID) (*User, error) {
	intUserID, err := rds.Get(UserClient{}.TableName(cid)).Int()
	if nil != err {
		return nil, err
	}

	userID := UserID(intUserID)
	user := &User{ID: userID}
	err = user.MakeUser()
	if nil != err {
		return nil, err
	}

	return user, nil
}

// UserClients 在线客户端连接
type UserClients []int32

// TableName 返回Redis key
func (UserClients) TableName() string {
	return "user:clients:userids"
}

// CreateUserClient 创建客户端用户关联关系
func CreateUserClient(uc UserClient) error {
	pipe := rds.Pipeline()
	defer pipe.Close()
	err := pipe.Set(uc.TableName(uc.ClientID), int32(uc.UserID), 0).Err()
	if nil != err {
		return err
	}
	pipe.SAdd(UserClients{}.TableName(), uc.ClientID)
	_, err = pipe.Exec()
	if nil != err {
		return err
	}

	return nil
}

// MultGetNames 返回当前登录的所有用户名
func MultGetNames() []string {
	ret, _ := rds.EvalSha(config.Config.Luas.MultGetNames, []string{UserClients{}.TableName()}, nil).Result()
	names, ins := []string{}, ret.([]interface{})

	for _, inter := range ins {
		names = append(names, inter.(string))
	}

	return names
}
