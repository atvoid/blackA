package main

import (
	"blackA/auth"
	"blackA/logging"
	"blackA/server"
	"fmt"
	"net"
)

var area string = "Main"

func main() {
	logging.StartLogging(logging.LOGLEVEL_DETAIL)
	listener, err := net.Listen("tcp", "10.0.1.4:789")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go server.GlobalRouter.StartRouting()
	fmt.Println("server started")
	for {
		conn, _ := listener.Accept()
		go auth.Authenticate(&conn)
	}
}
