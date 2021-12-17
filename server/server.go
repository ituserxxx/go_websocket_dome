package server

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	Conn      net.Conn
	Message   chan string   // 广播消息
	OnlineMap map[int]*User // 在线用户
	MapLock   sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[int]*User),
		Message: make(chan string),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Printf("%#v 111111111", err.Error())
		return
	}
	go func() {// 监听广播消息
		for {
			msg := <- s.Message
			s.MapLock.Lock()
			for _, u := range s.OnlineMap {
				u.Send(msg)
			}
			s.MapLock.Unlock()
		}
	}()
	for {//监听接收客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go s.handleConnect(conn)
	}
}
func (s *Server) handleConnect(conn net.Conn) {
	s.Conn = conn
	u := NewUser(rand.Intn(1000), conn, s)
	u.OnLine()
	s.Send(fmt.Sprintf("----【系统消息：欢迎%d进入房间成功】----",u.Id))
	//监听用户输入
	for {
		buf := make([]byte, 1024)
		n, err := u.Conn.Read(buf)
		if err != nil {
			u.OffLine()
			return
		}
		if n == 0 {
			u.OffLine()
			return
		}
		if err != nil && err != io.EOF {

			return
		}
		s.Message<- fmt.Sprintf("----用户  id=%d 用户说：【%s】",u.Id,string(buf[:n]))
	}
}


func (s *Server) Send(msg string) {
	_, err := s.Conn.Write([]byte("\n"+msg+"\n"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
