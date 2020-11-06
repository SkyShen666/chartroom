package process

import (
	"bufio"
	"client/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"shenguanchu.com/common/message"
)

func ShowMenu() {
	fmt.Printf("-------恭喜用户:%d 登录成功---------\n", CurUser.UserID)
	fmt.Println("-------1. 显示在线用户列表---------")
	fmt.Println("-------2. 发送消息---------")
	fmt.Println("-------3. 信息列表---------")
	fmt.Println("-------4. 退出系统---------")
	fmt.Println("请选择(1-4):")
	var choice int

	// 因为我们总会使用到SmsProcess实例,因此将其定义在switch外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &choice)
	switch choice {
	case 1:
		// 显示在线用户列表
		outputOnlineUser()
	case 2:
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Println("你想对大家说什么:)")
		content, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("inputReader.ReadString err =", err)
			return
		}
		// fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("消息列表")
	case 4:
		fmt.Println("你选择退出了系统...")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确")
	}
}

// serverProcessMes 和服务器保持通讯
func serverProcessMes(conn net.Conn) {
	// 创建一个Transfer实例,不停地读取服务器端发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在读取服务器端发送的消息...")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() err =", err)
			return
		}
		// 如果读取到消息,又是下一步处理逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType: // 有人上线了
			// 1. 取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			// 2. 把这个用户的信息,状态保存到客户map[int]User中
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType: // 有人群发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型...")
		}
	}
}
