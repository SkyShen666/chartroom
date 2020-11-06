package process2

import (
	"fmt"
)

var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 添加在线用户
func (um *UserMgr) AddOnlineUser(up *UserProcess) {
	um.onlineUsers[up.UserID] = up
}

// DelOnlineUser 删除某在线用户
func (um *UserMgr) DelOnlineUser(userID int) {
	delete(um.onlineUsers, userID)
}

// GetAllOnlineUsers 返回全部在线用户
func (um *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return um.onlineUsers
}

// GetOnlineByID 根据ID返回某在线用户
func (um *UserMgr) GetOnlineByID(userID int) (up *UserProcess, err error) {
	up, ok := um.onlineUsers[up.UserID]
	if !ok { // 不在线
		err = fmt.Errorf("用户%d 不在线", userID)
		return
	}
	return
}
