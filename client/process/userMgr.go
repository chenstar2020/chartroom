package process

import (
	"fmt"
	"gin_example/chartroom/client/model"
	"gin_example/chartroom/common/message"
)

var	onlineUsers map[int]*message.User
var curUser *model.CurUser

//初始化在线用户
func init(){
	onlineUsers = make(map[int] *message.User, 0)
	curUser = &model.CurUser{
		User: message.User{
			UserStatus: message.UserStatusOffline,  //初始化为离线状态
		},
	}
}

//添加在线用户
func addOnlineUser(newUser *message.User){
	onlineUsers[newUser.UserId] = newUser
}
//更新在线用户状态
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes){
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok{
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status

	onlineUsers[notifyUserStatusMes.UserId] = user
}

//显示在线用户
func displayOnlineUser(){
	fmt.Println("当前在线用户列表:")
	for id, user := range onlineUsers {
		fmt.Println("用户id:", id, "用户名称:", user.UserName, "状态:", user.UserStatus)
	}
}