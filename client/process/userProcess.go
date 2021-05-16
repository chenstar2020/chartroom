package process

import (
	"context"
	"encoding/json"
	"fmt"
	error2 "gin_example/chartroom/common/error"
	"gin_example/chartroom/common/message"
	"gin_example/chartroom/common/utils"
	"net"
)

type UserProcess struct {

}

func (this *UserProcess)Login(userId int, userPwd string)(err error){
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("[Login] net.Dial err=", err)
		return
	}
	defer conn.Close()

	tf := &utils.Transfer{
		Conn: conn,
	}

	//封装数据
	var mes message.Message
	mes.Type = message.LoginMesType

	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("[Login] loginMes json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("[Login] mes json.Marshal err=", err)
		return
	}

	err = tf.WritePkg(data)
	if err != nil {

	}
	fmt.Println("登录请求发送成功。。。")



	//处理服务器返回消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("[Login] tf.readPkg fail, err=", err)
		return
	}

	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("[Login] loginResMes json.Unmarshal fail, err=", err)
		return
	}

	if loginResMes.Code == error2.ERROR_OK ||
		loginResMes.Code == error2.ERROR_USER_LOGINED{  //已经登录过

		if loginResMes.Code == error2.ERROR_USER_LOGINED {
			fmt.Println("当前用户已经登录，其他设备将下线")
		}
		//保存当前用户  后续需要通过此用户发消息
		curUser.Conn = conn
		curUser.UserId = userId
		curUser.UserStatus = message.UserStatusOnline

		//保存在线用户列表
		for _, v := range loginResMes.UserIds {
			if v == userId {
				continue
			}
			user := &message.User{
				UserId: v,
				UserStatus: message.UserStatusOnline,
			}
			addOnlineUser(user)
		}

		quit := make(chan struct{})
		ctx, cancel := context.WithCancel(context.Background())

		//启动一个协程和和服务器保持通信
		go serverProcessMes(conn, quit)

		//启用一个协程显示主菜单
		go clientInputProcess(ctx, quit)

		<- quit
		cancel()    //此函数使得ShowMenu协程退出
	}else  {
		fmt.Println(loginResMes.Error)
	}

	return nil
}

func (this *UserProcess)Register(userId int, userName, userPwd string)(err error){
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()  //注册成功之后断开连接

	var mes message.Message
	mes.Type = message.RegisterMesType

	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserName = userName
	registerMes.User.UserPwd = userPwd

	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("Register json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("Register json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg fail, err=", err)
		return
	}
	fmt.Println("注册请求发送成功。。。")


	//处理服务器返回消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("Register readPkg err=", err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println("Register json.Unmarshal err=", err)
		return
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登录！")
	}else  {
		fmt.Println("注册失败:", registerResMes.Error)
	}
	return nil
}

//用户下线
func (this *UserProcess)Logout()(err error){
	curUser.UserStatus = message.UserStatusOffline

	var mes message.Message
	mes.Type = message.LogoutMesType

	var userLogoutMes message.UserLogoutMes
	userLogoutMes.User = curUser.User


	data, err := json.Marshal(userLogoutMes)
	if err != nil {
		fmt.Println("[Logout] userLogoutMes json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("[Logout] mes json.Marshal mes err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: curUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("[Logout] tf.WritePkg fail, err=", err)
		return
	}

	return
}
