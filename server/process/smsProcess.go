package process

import (
	"encoding/json"
	"fmt"
	"gin_example/chartroom/common/message"
	"gin_example/chartroom/common/utils"
	"net"
)

type SmsProcess struct {
	//暂时不需要字段
}

func (this *SmsProcess)SendGroupMes(mes *message.Message){
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId{
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess)SendMesToEachOnlineUser(data []byte, conn net.Conn){
	tf := utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
