package user

import (
	"blackA/logging"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

var area string = "User"

type User struct {
	Id           int
	Name         string
	conn         *net.Conn
	stopSigConn  chan bool
	stopSigMsg   chan bool
	userInput    chan Command
	UserInput    chan Command
	ServerInput  chan Command
	Disconnected bool
	RoomId       int
}

func CreateUser(id int, name string, conn *net.Conn) *User {
	u := User{Id: id, Name: name, conn: conn}
	u.ServerInput = make(chan Command, 10)
	u.userInput = make(chan Command, 10)
	u.stopSigConn = make(chan bool, 1)
	u.stopSigMsg = make(chan bool, 1)
	return &u
}

func (this *User) ResetConnection(conn *net.Conn) {
	this.stopSigConn <- true
	this.conn = conn
}

func (this *User) clearChannel() {
	for i := len(this.stopSigConn); i > 0; i-- {
		<-this.stopSigConn
	}
	for i := len(this.ServerInput); i > 0; i-- {
		<-this.ServerInput
	}
	for i := len(this.stopSigMsg); i > 0; i-- {
		<-this.stopSigMsg
	}
}

func (this *User) HandleConnection() {
	this.clearChannel()
	go this.receiveMsg()
	logging.LogInfo_Normal(area, fmt.Sprintf("Starting to handle user %v\n", this.Id))
PollLoop:
	for {
		select {
		case cc := <-this.userInput:
			// logging.LogInfo(area, fmt.Sprintf("Got Msg from %v with %v\n", this.Id, cc.ToMessage()))
			if cc.CmdType == CMDTYPE_PING {
				continue
			}
			this.UserInput <- cc
		case c := <-this.ServerInput:
			c.UserId = this.Id
			this.sendMsg(c)
		case <-this.stopSigConn:
			logging.LogInfo_Detail(area, fmt.Sprintf("Terminate handling connection from %v\n", this.Id))
			this.stopSigMsg <- true
			break PollLoop
		}
	}
	logging.LogInfo_Normal(area, fmt.Sprintf("End to handle user %v\n", this.Id))
}

func (this *User) receiveMsg() {
PollMsgLoop:
	for {
		select {
		case <-this.stopSigMsg:
			logging.LogInfo_Detail(area, fmt.Sprintf("Terminate receiving Msg from %v\n", this.Id))
			break PollMsgLoop
		default:
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
					logging.LogError(area, err.Error())
				}
			} else {
				logging.LogError(area, err.Error())
				this.userInput <- Command{CmdType: CMDTYPE_DISCONNECT, UserId: this.Id}
				this.stopSigConn <- true
				break PollMsgLoop
			}
		}
	}
}

func (this *User) sendMsg(c Command) {
	(*this.conn).Write([]byte(c.ToMessage() + "\u0001"))
}
