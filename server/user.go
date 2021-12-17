package server

import (
	"fmt"
	"net"
)

type User struct {
	Id     int
	Conn   net.Conn
	Server *Server
}

func NewUser(id int, c net.Conn, s *Server) *User {
	return &User{
		Id:     id,
		Conn:   c,
		Server: s,
	}
}

func (u *User) OnLine() {
	//锁相关参考：https://www.jianshu.com/p/679041bdaa39
	u.Server.MapLock.Lock()
	u.Server.OnlineMap[u.Id] = u
	u.Server.MapLock.Unlock()
	u.Server.Message <- fmt.Sprintf("--【用户%d上线了】",u.Id)
}
func (u *User) OffLine() {
	u.Server.MapLock.Lock()
	delete(u.Server.OnlineMap, u.Id)
	u.Server.MapLock.Unlock()
	u.Server.Message <- fmt.Sprintf("--【用户%d下线了】",u.Id)
}

func (u *User) Send(msg string) {
	_, err := u.Conn.Write([]byte("\n"+msg+"\n"))
	if err != nil {
		fmt.Println(err.Error())
	}
}

