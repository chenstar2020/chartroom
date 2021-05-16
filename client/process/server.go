package process

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_example/chartroom/common/message"
	"gin_example/chartroom/common/utils"
	"net"
	"sync"
)


//主菜单页面
func ShowMenu(){
	var key int
	var content string
	loop := true
	for{
		fmt.Println("*****主菜单******")
		fmt.Println("1. 在线用户列表")
		fmt.Println("2. 发送消息")
		fmt.Println("3. 信息列表")
		fmt.Println("4. 退出登录")
		fmt.Printf("请选择(1-4):")
		//发送消息是一个频繁操作，在外边创建实例
		smsProcess := &SmsProcess{}
		fmt.Scanln(&key)
		switch key {
		case 1:
			displayOnlineUser()
		case 2:
			fmt.Printf("请输入消息内容:")
			fmt.Scanln(&content)
			smsProcess.SendGroupMes(content)
		case 3:
			fmt.Println("信息列表")
		case 4:
			up := UserProcess{}
			up.Logout()
			loop = false
		default:
			fmt.Println("你输入的选项不正确")
		}

		if !loop {
			break
		}
	}
}

//和用户交互
func clientInputProcess(ctx context.Context, quit chan struct{}){
	var once sync.Once
	for {
		select {
		case <-ctx.Done():
			return
		default:
			go once.Do(func(){
				var key int
				var content string
				loop := true
				for{
					fmt.Println("*****主菜单******")
					fmt.Println("1. 在线用户列表")
					fmt.Println("2. 发送消息")
					fmt.Println("3. 信息列表")
					fmt.Println("4. 退出登录")
					fmt.Printf("请选择(1-4):")
					//发送消息是一个频繁操作，在外边创建实例
					smsProcess := &SmsProcess{}
					fmt.Scanln(&key)
					switch key {
					case 1:
						displayOnlineUser()
					case 2:
						fmt.Printf("请输入消息内容:")
						fmt.Scanln(&content)
						smsProcess.SendGroupMes(content)
					case 3:
						fmt.Println("信息列表")
					case 4:
						up := UserProcess{}
						up.Logout()
						loop = false
					default:
						fmt.Println("你输入的选项不正确")
					}

					if !loop {
						break
					}
				}

				quit<- struct{}{}
			})
		}
	}
}

//和服务器端保持通信
func serverProcessMes(conn net.Conn, quit chan struct{}){
	defer func(){
		fmt.Println("server process exit")
	}()
	//创建一个transfer实例  循环读取服务器消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for{
		mes, err := tf.ReadPkg()
		if err != nil {
			//fmt.Println("[serverProcessMes] tf.ReadPkg, err=", err)
			return
		}


		switch mes.Type {
		case message.NotifyUserStatusMesType:
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			displayGroupMes(&mes)
		case message.NotifyUserOfflineMesType:
			fmt.Println("收到下线消息...")
			curUser.UserStatus = message.UserStatusOffline
			quit<- struct{}{}
			return
		default:
			fmt.Println("服务器返回未知消息类型")
		}
	}
}
