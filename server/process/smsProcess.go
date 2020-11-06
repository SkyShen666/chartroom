package process2

import (
	"encoding/json"
	"fmt"
	"net"
	"server/utils"

	"shenguanchu.com/common/message"
)

type SmsProcess struct {
	// [暂时不需要字段]
}

// SendGroupMes 转发群聊消息
func (smsp *SmsProcess) SendGroupMes(mes *message.Message) {
	// 遍历服务器端的onlineUsers map[int]*UserProcess,
	// 将消息转发取出.
	// 取出mes的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserID {
			continue
		}
		smsp.SendMesToEachOnlineUser(data, up.Conn)
	}
}

// SendMesToEachOnlineUser 发送消息给其他所有在线用户
func (smsp *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个Transfer 实例，发送data
	tf := &utils.Transfer{
		Conn: conn, //
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
