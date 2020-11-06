package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// 定义几个用户状态常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

// Message 消息对象
type Message struct {
	Type string `json:"type"` // 消息的类型
	Data string `json:"data"` // 传输的数据
}

type LoginMes struct {
	UserID   int    `json:"userID"`   // 用户id
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

type LoginResMes struct {
	Code    int    `json:"code"` // 返回状态码,500:用户未注册,200:注册成功
	UsersID []int  // 保存用户id的切片
	Error   string `json:"error"` // 返回错误消息
}

type RegisterMes struct {
	User User `json:"user"` // 类型就是User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回状态码, 400:该用户已注册过, 200:注册成功
	Error string `json:"error"` // 返回错误信息
}

// NotifyUserStatusMes 配合服务器端推送用户状态变化的信息
type NotifyUserStatusMes struct {
	UserID int `json:"userID"`
	Status int `json:"status"`
}

// SmsMes 发送的消息
type SmsMes struct {
	Content string `json:"content"` // 内容
	User           // 匿名结构体,继承
}
