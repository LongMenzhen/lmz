package model

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v4"
)

// UserID 用户id
type UserID int32

// users 存储所有注册用户的id [1, 2, 3, 4, 5, 6]
// user:1 存储具体用户的基本信息;对应结构体 User

// NewUserID 返回自增的最新UserID
func NewUserID() UserID {
	key := "userid_seq"
	number, _ := rds.Incr(key).Result()

	return UserID(number)
}

// User 登陆用户
type User struct {
	ID       UserID `json:"id" msgpack:"id"`
	Username string `json:"username" msgpack:"username"`
	Password string `json:"-" msgpack:"password"`
}

// TableName 返回对应Redis存储的Key
func (User) TableName(id UserID) string {
	return fmt.Sprintf("user:%d", int(id))
}

// MakeUser 按照UserID创建用户
// 根据User.ID 查询Redis存储信息，将查询结果msgpack反解析为User 对象
func (u *User) MakeUser() error {
	bytea, _ := rds.Get(u.TableName(u.ID)).Bytes()
	if err := msgpack.Unmarshal(bytea, u); nil != err {
		return err
	}

	return nil
}

// MsgPackMarshal 序列化
func (u *User) MsgPackMarshal() ([]byte, error) {
	return msgpack.Marshal(u)
}

// CreateUser 创建用户
func CreateUser(user User) error {
	bytea, err := msgpack.Marshal(&user)
	if nil != err {
		return err
	}

	pipe := rds.Pipeline()
	defer pipe.Close()
	err = pipe.Set(user.TableName(user.ID), bytea, 0).Err()
	if nil != err {
		return err
	}
	err = pipe.SAdd(UserIDs{}.TableName(), int(user.ID)).Err()
	if nil != err {
		return err
	}
	_, err = pipe.Exec()
	if nil != err {
		return err
	}

	return nil
}

// NewUser 创建新的用户
func NewUser(username, password string) *User {
	return &User{
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
	return "userids"
}
