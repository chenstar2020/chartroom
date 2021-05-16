package main

import (
	"fmt"
	"gin_example/chartroom/common/message"
	"gin_example/chartroom/common/utils"
	process2 "gin_example/chartroom/server/process"
	"io"
	"net"
)

type Processor struct{
	Conn net.Conn
}

func (this *Processor)serverProcessMes( mes *message.Message)(err error){

	switch mes.Type{
	case message.LoginMesType:        //登录消息
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:     //注册消息
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:         //聊天消息
		smsProcess := process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	case message.LogoutMesType:      //用户下线消息
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		up.ServerProcessLogout(mes)
	default:
		fmt.Println("消息类型不存在，无法处理, MesType:", mes.Type)
	}
	return
}

func (this *Processor) process2()(err error){
	for {
		//将读取数据包封装成一个函数
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		var mes message.Message
		mes, err = tf.ReadPkg()
		if err != nil {         //一直处理数据 直到客户端主动断开连接
			if err == io.EOF {
				fmt.Println("客户端关闭了连接")
				return
			}
			fmt.Println("[process2] conn.Read err=", err)
			return
		}

		err = this.serverProcessMes(&mes)
		if err != nil {
			fmt.Println("[process2] serverProcessMes err=", err)
			return
		}
	}
}