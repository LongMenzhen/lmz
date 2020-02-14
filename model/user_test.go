package model

import (
	"testing"
)

func TestMsgPackMarshal(t *testing.T) {
	user := NewUser("chen", "123456")
	bytea, err := user.MsgPackMarshal()
	if nil != err {
		t.Errorf("msgpack序列化错误%s", err.Error())
	}

	t.Log(string(bytea))
}
