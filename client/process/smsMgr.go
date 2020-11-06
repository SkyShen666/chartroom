package process

import (
	"encoding/json"
	"fmt"

	"shenguanchu.com/common/message"
)

func outputGroupMes(mes *message.Message) { // 此处的mes一定是SmsMes
	//显示即可
	//1. 反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}

	//显示信息
	info := fmt.Sprintf("用户id:\t%d:\t%s", smsMes.UserID, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
