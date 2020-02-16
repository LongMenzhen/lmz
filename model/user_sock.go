package model

import (
	"fmt"

	"github.com/cyrnicolase/lmz/config"
	"github.com/cyrnicolase/lmz/engine"
)

// UserSock socket客户端与用户关联关系
type UserSock struct {
	SockID engine.SockID
	UserID UserID
}

// TableName 返回Rediskey
// user:client:{user_id}
func (UserSock) TableName(clientid engine.SockID) string {
	return fmt.Sprintf("user:client:%d", clientid)
}

// NewUserSock 创建客户端、顾客之间关系
func NewUserSock(clientID engine.SockID, userID UserID) *UserSock {
	return &UserSock{
		SockID: clientID,
		UserID: userID,
	}
}

// MakeUserBySockID 根据客户端id创建用户对象
func MakeUserBySockID(cid engine.SockID) (*User, error) {
	intUserID, err := rds.Get(UserSock{}.TableName(cid)).Int()
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

// UserSocks 在线客户端连接
type UserSocks []int32

// TableName 返回Redis key
func (UserSocks) TableName() string {
	return "user:clients:userids"
}

// CreateUserSock 创建客户端用户关联关系
func CreateUserSock(uc UserSock) error {
	pipe := rds.Pipeline()
	defer pipe.Close()
	err := pipe.Set(uc.TableName(uc.SockID), int32(uc.UserID), 0).Err()
	if nil != err {
		return err
	}
	pipe.SAdd(UserSocks{}.TableName(), uc.SockID)
	_, err = pipe.Exec()
	if nil != err {
		return err
	}

	return nil
}

// MultGetNames 返回当前登录的所有用户名
func MultGetNames() []string {
	ret, _ := rds.EvalSha(config.Config.Luas.MultGetNames, []string{UserSocks{}.TableName()}, nil).Result()
	names, ins := []string{}, ret.([]interface{})

	for _, inter := range ins {
		names = append(names, inter.(string))
	}

	return names
}
