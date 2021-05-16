package message

//定义user结构体
type User struct {
	UserId  int      `json:"userId"`   //用户id
	UserName string  `json:"userName"` //用户名称
	UserPwd string   `json:"userPwd"`  //用户密码
	UserStatus int   `json:"userStatus"`//用户状态
}

//用户状态
const (
	UserStatusOnline   = 1    //在线
	UserStatusOffline  = 2    //离线
	UserStatusBusy     = 3    //忙碌
)

//通知消息
const(
	NotifyMessageUserLoginTwice = "用户重复登录"
)