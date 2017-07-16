package user

import (
	"net"
	"bufio"
	"encoding/json"
	"fmt"
)

type User struct {
	Id			int
	Name		string
	conn		*net.Conn
	userInput	chan Command
	UserInput	chan Command
	ServerInput	chan Command
	Terminate	chan bool
}

func CreateUser(id int, name string, conn *net.Conn) User {
	u := User{ Id: id, Name: name, conn: conn }
	u.ServerInput = make(chan Command, 10)
	return u
}

func (this *User) HandleConnection() {
	go this.receiveMsg()
	select {
		case c := <- this.ServerInput:
			this.sendMsg(c)
		case c := <- this.userInput:
			fmt.Printf("Got Msg from %v\n",this.Id)
			this.UserInput <- c
		case <- this.Terminate:
			break
	}
}

func (this *User) receiveMsg() {
	for {
		msg, err := bufio.NewReader(*this.conn).ReadBytes('\u0001')
		fmt.Printf("Got Msg %v \n", string(msg))
		if err == nil {
			msg = msg[:len(msg)-1]
			//fmt.Printf("Got Msg %v \n", string(msg))
			var cmd Command
			err := json.Unmarshal(msg, &cmd)
			//fmt.Printf("%v", cmd.ToMessage())
			this.userInput <- cmd
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

func (this *User) sendMsg(c Command) {
	(*this.conn).Write([]byte(c.ToMessage()))
}