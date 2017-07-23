package main

import (
	"fmt"
	"net"
	"blackA/server"
	"blackA/user"
	"blackA/logging"
	"bufio"
	"encoding/json"
)

var area string = "Main"

func WaitingLogin(conn *net.Conn) {
	logging.LogInfo(area, fmt.Sprintf("client %v connected.\n", (*conn).RemoteAddr()))
	for {
		msg, err := bufio.NewReader(*conn).ReadBytes('\u0001')
		if err == nil {
			msg = msg[:len(msg)-1]
			var cmd user.Command
			err := json.Unmarshal(msg, &cmd)
			if err != nil {
				logging.LogError(area, err.Error())
				break
			} else if cmd.UserId == 0 {
				logging.LogInfo(area, fmt.Sprintf("invalid user id"))
				break
			} else {
				if cmd.CmdType == user.CMDTYPE_LOGIN {
					server.GlobalRouter.AddUser(cmd.UserId, "", conn)
					logging.LogInfo(area, fmt.Sprintf("Added User %v", cmd.UserId))
					return
				} else {
					break
				}
			}
		} else {
			logging.LogError(area, (err.Error()))
			break
		}
	}	
	(*conn).Close()
	logging.LogInfo(area, fmt.Sprintf("client %v Disconnected.\n", (*conn).RemoteAddr()))
}

func main() {
	logging.StartLogging()
	listener, err := net.Listen("tcp", "10.0.1.4:789")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go server.GlobalRouter.StartRouting()
	fmt.Println("server started")
	for {
		conn, _ := listener.Accept()
		go WaitingLogin(&conn)
	}
}
