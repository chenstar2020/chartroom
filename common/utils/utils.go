package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gin_example/chartroom/common/message"
	"net"
)

//将这些方法关联到结构体中

type Transfer struct {
	Conn net.Conn
	Buf [8096]byte   //数据缓冲
}

func (this *Transfer)ReadPkg()(mes message.Message, err error){

	n, err := this.Conn.Read(this.Buf[0:4])

	if err != nil {
		//fmt.Println("[ReadPkg] Conn.Read err=", err)
		return
	}

	//根据buf[:4] 转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkg读取消息内容
	n, err = this.Conn.Read(this.Buf[:pkgLen])
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("[ReadPkg] Conn.Read length error or ", err)
		err = errors.New("[ReadPkg] conn.Read read body error")
		return
	}

	//把pkglen 反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("[ReadPkg] mes json.Unmarshal err=", err)
		return
	}
	return
}

func (this *Transfer)WritePkg(data []byte)(err error){
	//1.先发送长度
	var pkgLen uint32
	pkgLen = uint32(len(data))

	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("[WritePkg] pkgLen conn.Write err:", err)
		return
	}

	//2.再发送data
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("[WritePkg]data conn.Write err:", err)
		return
	}

	return nil
}