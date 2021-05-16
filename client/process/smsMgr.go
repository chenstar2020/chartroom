package process

import (
	"encoding/json"
	"fmt"
	"gin_example/chartroom/common/message"
)

func displayGroupMes(mes *message.Message){
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("displayGroupMes json.Unmarshal err=", err)
		return
	}

	fmt.Printf("用户id:%v, 内容:%v\n", smsMes.UserId, smsMes.Content)
}
