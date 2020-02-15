package model

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v4"
)

func newGroupID() int32 {
	key := "groupid_seq"
	id := rds.Incr(key).Val()

	return int32(id)
}

// Group 组
type Group struct {
	ID int32 `json:"id"`
}

// TableName 返回记录Redis 的key
func (Group) TableName(id int32) string {
	return fmt.Sprintf("group:%d", int(id))
}

// CreateGroup 创建组
func CreateGroup(group Group) error {
	bytea, err := msgpack.Marshal(group)
	if nil != err {
		return err
	}

	pipe := rds.Pipeline()
	defer pipe.Close()
	err = pipe.Set(group.TableName(group.ID), bytea, 0).Err()
	if nil != err {
		return err
	}
	err = pipe.SAdd(GroupIDs{}.TableName(), int(group.ID)).Err()
	if nil != err {
		return err
	}
	_, err = pipe.Exec()
	if nil != err {
		return err
	}

	return nil
}

// NewGroup 创建新的组对象
func NewGroup() *Group {
	return &Group{
		ID: newGroupID(),
	}
}

// GroupIDs 返回组集合
type GroupIDs []int32

// TableName 返回组集合对应Redis key
func (GroupIDs) TableName() string {
	return "groupids"
}
