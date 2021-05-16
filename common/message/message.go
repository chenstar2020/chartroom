package message

import (
	error2 "gin_example/chartroom/common/error"
)

const (
	LoginMesType = "LoginMes"           //登录消息
	LoginResMesType = "LoginResMes"     //登录消息返回
	RegisterMesType = "RegisterMes"     //注册消息
	RegisterMesResType = "RegisterResMes"	//注册消息返回
	NotifyUserStatusMesType = "NotifyUserStatusMes"  //通知用户状态变更消息
	NotifyUserOfflineMesType = "NotifyUserOffline"   //通知用户下线消息
	SmsMesType = "SmsMes"   //聊天消息
	LogoutMesType = "LogoutMes"   //登出消息
)



//总的消息结构
type Message struct {
	Type string   	`json:"type"`//消息类型
	Data string     `json:"data"`//数据
}

//登录消息
type LoginMes struct {
	UserId int   	`json:"userId"`   //用户id
	UserPwd string `json:"userPwd"`   //用户密码
	UserName string `json:"userName"` //用户名称
}

//登录消息返回
type LoginResMes struct {
	Code    error2.Errno `json:"code"` //返回状态码
	Error   string       `json:"error"`
	UserIds []int        `json:"users"` //当前在线用户列表
}

//注册消息
type RegisterMes struct {
	User  User  `json:"user"`
}

//注册消息返回
type RegisterResMes struct {
	Code  error2.Errno `json:"code"`
	Error string       `json:"error"`
}

//用户状态变化消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

//聊天消息
type SmsMes struct {
	Content string `json:"content"`
	User   //只需要用到其中一些字段  定义成匿名结构体更方便
}

//通知用户下线消息
type NotifyUserLogoutMes struct {
	Reason string `json:"reason"`  //下线原因
}

//用户主动下线消息
type UserLogoutMes struct {
	Reason string `json:"reason"`
	User
}