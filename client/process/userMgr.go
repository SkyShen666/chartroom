package process

import (
	"client/model"
	"fmt"

	"shenguanchu.com/common/message"
)

// 客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser

// outputOnlineUser 在客户端显示当前在线的用户
func outputOnlineUser() {
	fmt.Println("当前在线用户列表如下:")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

// updateUserStatus 处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserID]
	if !ok { // 原来没有
		user = &message.User{
			UserID: notifyUserStatusMes.UserID,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserID] = user

	outputOnlineUser()
}
