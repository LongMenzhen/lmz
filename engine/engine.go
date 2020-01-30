package engine

// HandlerFunc 执行函数
type HandlerFunc func(ctx *Context)

// Action 操作
type Action struct {
	Event string
	F     HandlerFunc
}

// Actions 注册执行方法集合
var Actions []Action

// Registe 注册
func Registe(event string, f HandlerFunc) {
	Actions = append(Actions, Action{
		Event: event,
		F:     f,
	})
}
