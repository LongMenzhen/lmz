package app

import (
	"net/http"

	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// Register 注册用户
func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	user := model.NewUser()
	if err := user.Save(); nil != err {
		logrus.Error("注册用户信息保存失败")
		return
	}

}
