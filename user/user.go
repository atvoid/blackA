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
	Disconnected	bool
	RoomId		int
}

func CreateUser(id int, name string, conn *net.Conn) *User {
	u := User{ Id: id, Name: name, conn: conn }
	u.ServerInput = make(chan Command, 10)
	u.userInput = make(chan Command, 10)
	return &u
}

func (this *User) HandleConnection() {
	//userInput := make(chan Command, 10)
	go this.receiveMsg()
	fmt.Printf("Starting to handl user %v\n", this.Id)
	PollLoop:
	for {
		select {
			case cc := <- this.userInput:
				fmt.Printf("Got Msg with %v\n", cc.ToMessage())
				fmt.Printf("Got Msg from %v\n",this.Id)
				this.UserInput <- cc
			case c := <- this.ServerInput:
				c.UserId = this.Id
				this.sendMsg(c)
			case <- this.Terminate:
				fmt.Printf("Terminate Msg from %v\n",this.Id)
				break PollLoop
		}
	}
	fmt.Println("End")
}

func (this *User) receiveMsg() {
	for {
		msg, err := bufio.NewReader(*this.conn).ReadBytes('\u0001')
		//fmt.Printf("Got Msg %v \n", string(msg))
		if err == nil {
			msg = msg[:len(msg)-1]
			//fmt.Printf("Got Msg %v \n", string(msg))
			var cmd Command
			err := json.Unmarshal(msg, &cmd)
			//fmt.Printf("%v", cmd.ToMessage())
			//uInput <- cmd
			cmd.UserId = this.Id
			this.userInput <- cmd
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
			this.userInput <- Command{ CmdType: CMDTYPE_DISCONNECT, UserId: this.Id }
			break
		}
	}
}

func (this *User) sendMsg(c Command) {
	(*this.conn).Write([]byte(c.ToMessage() + "\u0001"))
}