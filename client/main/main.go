package main

import (
	"fmt"
	"gin_example/chartroom/client/process"
)

func main(){
	//一系列初始化

	var key int
	var loop = true
	for {
		fmt.Println("******在线聊天系统******")
		fmt.Println("		1.登录")
		fmt.Println("		2.注册")
		fmt.Println("		3.退出")
		fmt.Printf("请选择(1-3):")
		var userId int
		var userPwd string
		var userName string
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Printf("请输入用户id:")
			fmt.Scanln( &userId)
			fmt.Printf("请输入用户密码:")
			fmt.Scanln(&userPwd)
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Printf("请输入用户id:")
			fmt.Scanln( &userId)
			fmt.Printf("请输入用户名字:")
			fmt.Scanln(&userName)
			fmt.Printf("请输入用户密码:")
			fmt.Scanln(&userPwd)
			up := &process.UserProcess{}
			up.Register(userId, userName, userPwd)
		case 3:
			fmt.Printf("退出操作\n")
			loop = false
		default:
			fmt.Printf("输入有误，请重新输入\n")
		}


		if !loop {
			break
		}
	}


}