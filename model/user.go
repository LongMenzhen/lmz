package model

// User 登陆用户
type User struct {
	Name     string `json:"name"`
	ClientID int32  `json:"client_id"`
}
