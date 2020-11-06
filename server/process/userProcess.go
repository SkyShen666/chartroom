// Package process2 原因: process会和process()函数重名
package process2

import (
	"encoding/json"
	"fmt"
	"net"
	"server/model"
	"server/utils"

	"shenguanchu.com/common/message"
)

type UserProcess struct {
	Conn net.Conn
	// 表示Conn是哪个用户
	UserID int
}

// ServerProcessLogin :处理登录消息
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 核心代码
	// 1. 先从mes中取出mes.Data,并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return
	}

	// 2. 校验登录信息,并初始化返回消息
	// 2.1 先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginMesType
	// 2.2 声明一个LoginResMes,并完成赋值
	var loginResMes message.LoginResMes

	// 2.3 校验密码
	// 到redis中去完成验证
	// 使用model.MyUserDao到redis中验证
	user, err := model.MyUserDao.Login(loginMes.UserID, loginMes.UserPwd)
	if err != nil {
		switch err {
		case model.ERROR_USER_NOTEXISTS:
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		case model.ERROR_USER_PWD:
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		default:
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
	} else {
		loginResMes.Code = 200
		// up.UserID = loginMes.UserID
		fmt.Println(user, "登录成功!")
		// 登录成功,要把该登录成功的用户放入userMgr中
		// 将登录成功的用户的userID 赋给 up
		up.UserID = loginMes.UserID
		userMgr.AddOnlineUser(up)

		// 通知其他的在线用户,我上线了
		up.NotifyOthersOnlineUsers(loginMes.UserID)

		// 将当前在线的用户的id,放入到loginResMes.UsersID中
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersID = append(loginResMes.UsersID, id)
		}
	}

	// 3. 将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 4. 将data赋值给resMes
	resMes.Data = string(data)

	// 5. 对resMes进行序列化,准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 6. 发送data,我们将其封装到writePkg函数
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		return
	}
	return
}

func (up *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1. 先从mes中取出mes.Data,并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err =", err)
		return
	}

	// 2. 声明resMes,registerResMes(下面会用到)
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	// 3. 到redis中完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		// 注册成功
		registerResMes.Code = 200
		fmt.Println("注册成功...")
	}

	// 4. 向服务器发送注册后的返回信息
	// 4.1 序列化registerResMes
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err =", err)
		return
	}

	// 4.2 将data赋值给resMes
	resMes.Data = string(data)

	// 4.3 对resMes序列化,准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err =", err)
		return
	}

	// 4.4 发送data,将其封装到writePkg函数
	// 因为使用mvc分层模式,先创建一个Transfer实例,然后读取
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg err =", err)
		return
	}
	return
}

// NotifyOthersOnlineUsers 通知所有在线用户
// userID 通知其他用户,此用户上线
func (up *UserProcess) NotifyOthersOnlineUsers(userID int) {
	// 遍历onlineUsers,然后一个个发送NotifyUsersStatusMes
	for id, up := range userMgr.onlineUsers {
		if id == userID {
			continue
		}
		// 开始通知
		up.NotifyMeOnline(userID)
	}
}

// NotifyMeOnline 通知其他人,'我'上线了
func (up *UserProcess) NotifyMeOnline(userID int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userID
	notifyUserStatusMes.Status = message.UserOnline

	// 将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}
	// 将序列化后的nitifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)

	// 对mes再次进行序列化,准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	// 发送,创建Transfer实例
	tf := &utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err =", err)
		return
	}
}
