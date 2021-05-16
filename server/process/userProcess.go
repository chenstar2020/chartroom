package process

import (
	"encoding/json"
	"fmt"
	error2 "gin_example/chartroom/common/error"
	"gin_example/chartroom/common/message"
	"gin_example/chartroom/common/utils"
	"gin_example/chartroom/server/model"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	UserId int
}
//处理登录
func (this *UserProcess)ServerProcessLogin(mes *message.Message)(err error) {
	//1.取出data,并反序列化
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("[ServerProcessLogin] loginMes json.Unmarshal  fail, err=", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes

	//2.校验用户名和密码
	_, errno := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if errno != error2.ERROR_OK {
		loginResMes.Code = errno
		loginResMes.Error = error2.ErrorNoMsgMap[errno]
	}else{
		 onlineUser := userMgr.GetOnlineUserByUserId(loginMes.UserId)
		if onlineUser != nil {
			loginResMes.Code = error2.ERROR_USER_LOGINED
			loginResMes.Error = error2.ErrorNoMsgMap[error2.ERROR_USER_LOGINED]

			//删除已经登录的用户
			userMgr.DelOnlineUser(onlineUser.UserId)
			//通知已经登录的用户下线
			this.NotifyOnlineUserLogout(onlineUser)
			goto END
		}
		loginResMes.Code = error2.ERROR_OK
		loginResMes.Error = error2.ErrorNoMsgMap[error2.ERROR_OK]

		//添加在线用户
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)

		//通知其他用户自己上线
		this.NotifyOthersOnlineUser(this.UserId)

		//返回在线用户列表
		for k, _ := range userMgr.onlineUsers{
			loginResMes.UserIds = append(loginResMes.UserIds, k)
		}
	}
END:
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("[ServerProcessLogin] loginResMes json.Marshal fail ,err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("[ServerProcessLogin] resMes json.Marshal fail ,err=", err)
		return
	}

	//3.发送data
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//处理注册
func (this *UserProcess)ServerProcessRegister(mes *message.Message)(err error) {
	//1.取出data,并反序列化
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("[ServerProcessRegister] registerMes json.Unmarshal  fail, err=", err)
		return
	}
	fmt.Println("服务器收到用户注册消息，用户Id:", registerMes.User.UserId, "用户密码:", registerMes.User.UserPwd)

	var resMes message.Message
	resMes.Type = message.RegisterMesResType
	var registerResMes message.RegisterResMes

	//2.注册用户
	errno := model.MyUserDao.Register(&registerMes.User)
	fmt.Println("注册用户结果：", error2.ErrorNoMsgMap[errno])


	registerResMes.Code = errno
	registerResMes.Error = error2.ErrorNoMsgMap[errno]


	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("[ServerProcessRegister] registerResMes json.Marshal fail ,err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("[ServerProcessRegister] resMes json.Marshal fail ,err=", err)
		return
	}

	//3.发送data
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("[ServerProcessRegister] tf.WritePkg fail ,err=", err)
		registerResMes.Code = error2.ERROR_SEVER_ERR
		registerResMes.Error = error2.ErrorNoMsgMap[error2.ERROR_SEVER_ERR]
		return
	}
	return
}

//用户下线
func (this *UserProcess) ServerProcessLogout(mes *message.Message)(err error){
	//1.取出data,并反序列化
	var userLogoutMes message.UserLogoutMes
	err = json.Unmarshal([]byte(mes.Data), &userLogoutMes)
	if err != nil {
		fmt.Println("[ServerProcessLogout] userLogoutMes json.Unmarshal fail, err=", err)
		return
	}
	fmt.Println("服务器收到用户下线消息，用户ID:", userLogoutMes.UserId)

	//删除已经登录的用户
	userMgr.DelOnlineUser(userLogoutMes.UserId)
	return
}

//通知所有在线用户
func (this *UserProcess) NotifyOthersOnlineUser(userId int){
	for id, up := range userMgr.onlineUsers{
		if id == userId{
			continue
		}

		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int){
	//组装消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserStatusOnline

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
	}
	return
}

//通知用户下线
func (this *UserProcess) NotifyOnlineUserLogout(user *UserProcess){
	var mes message.Message
	mes.Type = message.NotifyUserOfflineMesType

	var notifyUserLogoutMes message.NotifyUserLogoutMes
	notifyUserLogoutMes.Reason = message.NotifyMessageUserLoginTwice

	data, err := json.Marshal(notifyUserLogoutMes)
	if err != nil {
		fmt.Println("[NotifyOnlineUserLogout] notifyUserLogoutMes json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("[NotifyOnlineUserLogout] mes json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: user.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("[NotifyOnlineUserLogout] WritePkg err=", err)
	}
	return
}

