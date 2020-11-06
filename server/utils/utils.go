package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"shenguanchu.com/common/message"
)

// Transfer 将传输数据的方法关联到此对象中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte
}

// WritePkg 传输数据(data)
func (tf *Transfer) WritePkg(data []byte) (err error) {
	// 1. 先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))

	binary.BigEndian.PutUint32(tf.Buf[:4], pkgLen)
	// 2. 发送长度
	n, err := tf.Conn.Write(tf.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	// 3. 发送消息(data)本身
	n, err = tf.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	return
}

// ReadPkg 读取数据
func (tf *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据...")
	// conn.Read 在conn没有被关闭的情况下,才会阻塞
	// 如果客户端关闭了conn,就不会阻塞

	// 1. 读取数据长度(用于校验是否丢包)
	_, err = tf.Conn.Read(tf.Buf[:4])
	if err != nil {
		// fmt.Println("conn.Read err =", err)
		return
	}

	// 2. 计算消息长度,根据buf[:4]转成一个unit32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(tf.Buf[:4])

	// 3. 根据pkgLen读取消息内容
	n, err := tf.Conn.Read(tf.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		err = errors.New("read pkg body error")
		return
	}

	// 4. 把pkgLen反序列化成-> message.Message
	// 注意!!!!:&mes
	err = json.Unmarshal(tf.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return
	}
	return
}
