package main

import (
	"fmt"
	"net"
	"blackA/server"
	"blackA/user"
	"bufio"
	"encoding/json"
)

func WaitingLogin(conn *net.Conn) {
	fmt.Printf("client %v connected.\n", (*conn).RemoteAddr())
	for {
		msg, err := bufio.NewReader(*conn).ReadBytes('\u0001')
		if err == nil {
			msg = msg[:len(msg)-1]
			var cmd user.Command
			err := json.Unmarshal(msg, &cmd)
			if err != nil {
				fmt.Println(err.Error())
				break
			} else if cmd.UserId == 0 {
				fmt.Println("invalid user id")
				break
			} else {
				if cmd.CmdType == user.CMDTYPE_LOGIN {
					server.GlobalRouter.AddUser(cmd.UserId, "", conn)
					fmt.Println("Added User")
					return
				} else {
					break
				}
			}
		} else {
			fmt.Println(err.Error())
			break
		}
	}	
	(*conn).Close()
	fmt.Printf("client %v Disconnected.\n", (*conn).RemoteAddr())
}

func main() {
	listener, _ := net.Listen("tcp", "192.168.199.189:789")
	go server.GlobalRouter.StartRouting()
	for {
		conn, _ := listener.Accept()
		go WaitingLogin(&conn)
	}
}
