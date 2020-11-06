package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"shenguanchu.com/common/message"
)

// 我们服务器启动后,就初始化一个userDao实例
// 把它做成全局变量,在需要和redis操作时,直接使用即可
var (
	MyUserDao *UserDao
)

// UserDao 完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// NewUserDao 使用工厂模式,创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// getUserByID 根据用户id,返回一个user实例+err
func (ud *UserDao) getUserByID(conn redis.Conn, id int) (user *message.User, err error) {
	// 通过给定的id,从redis中查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		// redis.ErrNil 表示在users哈希表中,没有找到对应id
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &message.User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return
	}
	return
}

// Login 完成登录的校验
func (ud *UserDao) Login(userID int, userPwd string) (user *message.User, err error) {
	// 1. 先从UserDao的连接池中取出一个连接
	conn := ud.pool.Get()
	defer conn.Close()
	user, err = ud.getUserByID(conn, userID)
	if err != nil {
		return
	}
	// 2. 证明这个用户是获取到的
	// 2.1 如果用户的id和pwd都正确，则返回一个user实例
	// 2.2 如果用户的id或pwd有错误，则返回对应的错误信息
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (ud *UserDao) Register(user *message.User) (err error) {
	// 1. 从UserDao的连接池中取出一个连接
	conn := ud.pool.Get()
	defer conn.Close()

	// 2. 检查此用户是否已经注册过
	_, err = ud.getUserByID(conn, user.UserID)
	// 2.1 err == nil 说明查到了此用户,即已经注册过
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	// 3. 没有注册过,可以完成注册
	// 3.1 序列化
	//fmt.Println("userDao, 没有注册过")
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	// 3.2 入库
	_, err = conn.Do("HSet", "users", user.UserID, string(data))
	if err != nil {
		fmt.Println("保存注册用户信息错误 err =", err)
		return
	}
	return
}
