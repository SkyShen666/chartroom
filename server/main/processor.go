package main

import (
	"fmt"
	"io"
	"net"
	process2 "server/process"
	"server/utils"

	"shenguanchu.com/common/message"
)

type Processor struct {
	Conn net.Conn
}

// serverProcessMes: 根据客户端发送消息种类的不同,决定调用哪个函数来处理
func (p *Processor) serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		// 处理登录消息
		// 创建一个UserProcess实例
		up := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(mes)
		if err != nil {
			return
		}
	case message.RegisterMesType:
		// 处理注册消息
		up := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRegister(mes)
		if err != nil {
			return
		}
	case message.SmsMesType:
		// 创建一个SmsProcess实例完成转发群聊消息
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在,无法处理...")
	}
	return
}

// process2 总处理方法,处理客户端发送的消息
// 用process2原因: 同一包下(package main),main.go有process函数,这里不要重名
func (p *Processor) process2() (err error) {
	// 循环的处理客户端发送的消息
	for {
		// 注意tf应放在循环里面,要不然每次都是读的同一个conn的消息
		tf := &utils.Transfer{
			Conn: p.Conn,
		}
		// 这里将读取数据包,直接封装成一个函数readPkg(),返回Message, Err
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出,服务端也退出...")
				return err
			} else {
				fmt.Println("readPkg err =", err)
				return err
			}
		}
		//fmt.Println("mes =", mes)
		err = p.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
