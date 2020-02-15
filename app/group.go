package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// CreateGroup 创建组
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	if "POST" == r.Method {
		group := model.NewGroup()
		if err := model.CreateGroup(*group); nil != err {
			logrus.Error("创建消息组失败" + err.Error())
			fmt.Fprint(w, "注册用户失败")
			return
		}

		fmt.Fprint(w, "创建消息组成功: "+strconv.Itoa(int(group.ID)))
	}
}
