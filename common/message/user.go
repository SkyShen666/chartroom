package message

// User 定义一个用户对象
type User struct {
	// 为了序列化何反序列化成功，必须保证
	// 用户信息的json字符串的key 和 结构体的字段对应的 tag 名字一致
	UserID     int    `json:"userID"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"`
	Sex        string `json:"sex"`
}
