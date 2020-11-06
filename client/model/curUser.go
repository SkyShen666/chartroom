package model

import (
	"net"

	"shenguanchu.com/common/message"
)

// CurUser 因为在客户端,很多地方会使用到curUser,将其作为一个全局变量
type CurUser struct {
	Conn net.Conn
	message.User
}
