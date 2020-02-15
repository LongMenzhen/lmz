package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// Register 注册用户
func Register(w http.ResponseWriter, r *http.Request) {
	if "POST" == r.Method {
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		user := model.NewUser(username, password)
		if err := model.CreateUser(*user); nil != err {
			logrus.Error("注册用户信息保存失败" + err.Error())
			fmt.Fprint(w, "注册用户失败")
			return
		}

		ret, _ := json.Marshal(user)
		fmt.Fprint(w, string(ret))
	}
}
