package client

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	Ip   string
	Port int
	Conn net.Conn
}

func NewClient(ip string, port int) *Client {
	return &Client{
		Ip:   ip,
		Port: port,
	}
}

func (c *Client) Run() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Ip, c.Port))
	if err != nil {
		fmt.Printf("1111 %#v", err.Error())
		return
	}
	c.Conn = conn
	defer conn.Close()
	_, err = conn.Write([]byte("我驾着七彩祥云来了"))
	if err != nil {
		fmt.Printf("2222 %#v", err.Error())
		return
	}
	// 监听服务器发送的消息
	go func() {
		_, _ = io.Copy(os.Stdout, c.Conn)
	}()
	_, _ = c.Conn.Write([]byte("Succ"))
	for {
		if c.menu() {
			break
		}
	}
}
func (c *Client) menu() bool {
	fmt.Println(">>> 请输入内容：")
	var flag string
	_, _ = fmt.Scanln(&flag)
	if flag == "exit"{
		return true
	}
	_, _ = c.Conn.Write([]byte(flag))
	return false
}
