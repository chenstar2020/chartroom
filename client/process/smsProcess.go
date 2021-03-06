package process

import (
	"encoding/json"
	"fmt"
	"gin_example/chartroom/common/message"
	"gin_example/chartroom/common/utils"
)

type SmsProcess struct {

}

func (this *SmsProcess)SendGroupMes(content string)(err error){
	var mes message.Message
	mes.Type = message.SmsMesType

	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = curUser.UserId
	smsMes.UserStatus = curUser.UserStatus

	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal mes err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: curUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg fail, err=", err)
		return
	}
	return
}