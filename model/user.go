package model

import (
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack"
)

// users 存储所有注册用户的id [1, 2, 3, 4, 5, 6]
// user:1 存储具体用户的基本信息;对应结构体 User

// UserID 对应用户id
type UserID int32

// NewUserID 返回自增的最新UserID
func NewUserID() UserID {
	key := "userid"
	number, _ := rds.Incr(key).Result()

	return UserID(number)
}

// User 登陆用户
type User struct {
	ID        UserID    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 返回对应Redis存储的Key
func (User) TableName(id int32) string {
	return fmt.Sprintf("user:%d", int(id))
}

// Save 保存
func (u *User) Save() error {
	bytea, err := msgpack.Marshal(u)
	if nil != err {
		return err
	}

	err = rds.Set(u.TableName(u.ID), bytea, 0).Err()
	if nil != err {
		return err
	}
}

// NewUser 创建新的用户
func NewUser(username, password string) User {
	return User{
		ID:       NewUserID(),
		Username: username,
		Password: password,
	}
}

// UserIDs 用户集合
// 在Redis中存储为集合存储，存储所有的UserID
type UserIDs []UserID

// TableName 返回对应Redis存储的Key
func (UserIDs) TableName() string {
	return "users"
}
