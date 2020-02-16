package model

// ErrGroupNotFound 组不存在
type ErrGroupNotFound struct {}

// Error 返回错误信息
func (e ErrGroupNotFound) Error() string {
	return "组不存在"
}