package main

import (
	"fmt"
	"gin_example/chartroom/server/db"
	"gin_example/chartroom/server/model"
	"net"
	"os"
)

func process(conn net.Conn) {
	defer conn.Close()

	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯关闭, err=", err)
	}
}


func main(){

	//初始化配置

	//初始化数据库
	db.InitPool("localhost:6379", 16, 0, 300)
	model.MyUserDao = model.NewUserDao(db.RedisPool)

	fmt.Println("服务器在8889监听。。。")
	service := ":8889"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listen, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//启动一个协程,处理客户端连接
		go process(conn)
	}

}

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s", err.Error())
		os.Exit(1)
	}
}