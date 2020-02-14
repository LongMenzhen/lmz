package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyrnicolase/lmz/model"
	"github.com/sirupsen/logrus"
)

// Register 注册用户
func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	user := model.NewUser(username, password)
	if err := model.CreateUser(*user); nil != err {
		logrus.Error("注册用户信息保存失败" + err.Error())
		fmt.Fprint(w, "注册用户失败")
		return
	}

	nuser := &model.User{
		ID: user.ID,
	}
	nuser.MakeUser()
	log.Println("-->", nuser)

	ret, _ := json.Marshal(user)
	log.Println(string(ret))
	// w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, string(ret))
}
