package main

import (
	"fmt"
	"net"
	"server/model"
	"time"
)

// init 先于main启动,进行初始化操作
func init() {
	// 当服务器启动时,我们就去初始化我们的redis连接池
	initPool("127.0.0.1:6379", 16, 0, 300*time.Second)
	initUserDao()
}

// initUserDao UserDao初始化
func initUserDao() {
	// pool 为全局变量
	// 这里需注意初始化顺序问题
	// initPool, 在initUserDao之前
	model.MyUserDao = model.NewUserDao(pool)
}

func process(conn net.Conn) {
	// 延时关闭conn
	defer conn.Close()

	// 创建总控实例(Processor)
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器端通信协程错误, err =", err)
		return
	}
}

func main() {
	fmt.Println("服务器在8889端口监听...")
	listener, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listener.Close() 

	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}

	// 一旦监听成功,就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err =", err)
		}
		// 一旦连接成功,则启动一个协程和客户端保持通讯.
		go process(conn)
	}
}
