package main

import (
	"client/process"
	"fmt"
	"os"
)

// 用户id,密码和用户名
var userID int
var userPwd string
var userName string

func main() {
	// 接收用户的选择
	var choice int

	for true {
		fmt.Println("------------------欢迎登录多人聊天系统------------------")
		fmt.Println("\t\t\t 1 登录聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &choice)
		switch choice {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userID)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			// 完成登录
			up := &process.UserProcess{}
			up.Login(userID, userPwd)
		case 2:
			// 1. 输入注册信息
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userID)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名(nickname):")
			fmt.Scanf("%s\n", &userName)
			// 2. 调用UserProcess,完成注册请求
			up := &process.UserProcess{}
			up.Register(userID, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("您的输入有误,请重新输入")
		}
	}
}
