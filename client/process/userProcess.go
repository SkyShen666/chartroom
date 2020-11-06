package process

import (
	"client/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"shenguanchu.com/common/message"
)

type UserProcess struct {
	// 暂时不需要字段
}

// Login 给UserProcess对象,关联一个用户登录的方法
func (up *UserProcess) Login(userID int, userPwd string) (err error) {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err =", err)
		return
	}
	defer conn.Close()

	// 2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3. 创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	// 4. 将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 5.将序列化后的loginMes,即data赋给mes.Data字段
	mes.Data = string(data)

	// 6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 7. 创建Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	tf.WritePkg(data)

	// 这里还需处理服务器端返回的消息
	// 8. 读取服务器返回的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("ReadPkg err =", err)
		return
	}

	// 9. 将mes的Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")

		// 初始化CurUser
		CurUser.Conn = conn
		CurUser.UserID = userID
		CurUser.UserStatus = message.UserOnline

		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMes.UsersID {
			if v == userID {
				continue
			}
			fmt.Println("用户id:\t", v)
			// 完成客户端的 onlineUsers 初始化
			user := &message.User{
				UserID:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		// 这里需要在客户端启动一个协程
		// 该协程保持和服务器的通讯,如果服务器有数据推送给客户端
		// 则接受并显示在客户端的终端.
		go serverProcessMes(conn)

		// 显示登录成功的菜单[循环]
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

func (up *UserProcess) Register(userID int, userPwd string, userName string) (err error) {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err =", err)
		return
	}
	// 1.1 延时关闭
	defer conn.Close()

	// 2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType

	// 3. 创建一个RegisterMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. 将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 5. 把data赋给mes.Data字段
	mes.Data = string(data)

	// 6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 7. 创建Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 8. 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册信息发送错误 err =", err)
		return
	}

	// 9. 读取服务器返回的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("tf.ReadPkg err =", err)
		return
	}

	// 10. 将mes的Data部分反序列化成 RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功, 请登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
