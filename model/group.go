package model

import (
	"fmt"

	"github.com/cyrnicolase/lmz/config"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v4"
)

func newGroupID() GroupID {
	key := "groupid_seq"
	id := rds.Incr(key).Val()

	return GroupID(id)
}

// GroupID 组id
type GroupID int32

// Group 组
type Group struct {
	ID        GroupID `json:"id" msgpack:"id"`
	Name      string  `json:"name" msgpack:"name"`
	CreatorID UserID  `json:"creator_id" msgpack:"creator_id"`
}

// TableName 返回记录Redis 的key
func (Group) TableName(id GroupID) string {
	return fmt.Sprintf("group:%d", int(id))
}

// MakeGroup 根据id 查询Redis返回组结构体
func (g *Group) MakeGroup() error {
	msg, err := rds.Get(g.TableName(g.ID)).Bytes()
	if nil != err {
		return ErrGroupNotFound{}
	}

	err = msgpack.Unmarshal(msg, g)
	if nil != err {
		return err
	}

	return nil
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

// MultGroups 获取所有的组对象
func MultGroups() []Group {
	result, err := rds.EvalSha(config.Config.Luas.MultGroups, []string{GroupIDs{}.TableName()}, nil).Result()
	if nil != err {
		logrus.Error("读取Redis批量组信息失败" + err.Error())
		return nil
	}

	groups, ret := []Group{}, result.([]interface{})
	for _, msg := range ret {
		data := msg.(string)
		var group Group
		if err := msgpack.Unmarshal([]byte(data), &group); nil != err {
			continue
		}

		groups = append(groups, group)
	}

	return groups
}

// NewGroup 创建新的组对象
func NewGroup(user *User, name string) *Group {
	return &Group{
		ID:        newGroupID(),
		Name:      name,
		CreatorID: user.ID,
	}
}

// GroupIDs 返回组集合
type GroupIDs []GroupID

// TableName 返回组集合对应Redis key
func (GroupIDs) TableName() string {
	return "groupids"
}

// GroupUsers 组下关联用户
type GroupUsers []UserID

// TableName 返回具体组下的用户集合key
func (GroupUsers) TableName(groupID GroupID) string {
	return fmt.Sprintf("group:users:%d", int(groupID))
}

// AddUserToGroup 将用户加入到组
func AddUserToGroup(user *User, group *Group) error {
	key := GroupUsers{}.TableName(group.ID)
	err := rds.SAdd(key, user.ID).Err()
	if nil != err {
		return err
	}

	return nil
}
