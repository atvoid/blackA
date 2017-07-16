package room

import (
	"blackA/user"
	"blackA/game"
	"net"
	"encoding/json"
	"fmt"
)

type Room struct {
	Id			int
	Game		game.CardGame
	Users		[]user.User
	msgReceiver	chan user.Command
}

func CreateRoom() Room {
	room := Room{ msgReceiver: make(chan user.Command, 20)}
	room.Game = game.CreateCardGame(0)
	room.Users = make([]user.User, 0, len(room.Game.Players))
	return room
}

func (this *Room) IsFull() bool {
	fmt.Printf("Users: %v, Cap: %v\n", len(this.Users), cap(this.Users))
	if len(this.Users) < cap(this.Users) {
		return false
	}
	return true
}

func (this *Room) AddUser(id int, name string, conn *net.Conn) bool {
	if (!this.IsFull()) {
		this.Users = append(this.Users, user.CreateUser(id, name, conn))
		u := &this.Users[len(this.Users)-1]
		u.UserInput = this.msgReceiver
		//go u.ReceiveMsg()
		go u.HandleConnection()
		return true
	} else {
		return false
	}
}

func (this *Room) getUserIndex(id int) int {
	for i, v := range this.Users {
		if v.Id == id {
			return i
		}
	}
	panic("No Such User")
}

func (this *Room) notifyUsers(cmd user.Command) {
	for _, v := range this.Users {
		v.ServerInput <- cmd
		//fmt.Println(cmd.ToMessage())
	}
}

func (this *Room) notifyStatus() {
	for i, v := range this.Users {
		msg := this.Game.GetStatus(i)
		v.ServerInput <- user.Command{ Id: 0, Command: msg.ToMessage() }
		// fmt.Println(msg.ToMessage())
	}
}

func (this *Room) Start() {
	this.Game.Start()
	this.notifyStatus()
	PollLoop:
	for {
		select {
			case c:= <- this.msgReceiver:
				fmt.Println(c.ToMessage())
				var cmd game.CardCommand
				idx := this.getUserIndex(c.UserId)
				json.Unmarshal([]byte(c.Command), &cmd)
				if cmd.CmdType == game.CMDTYPE_PASS {
					result := this.Game.Pass(idx)
					if (result) {
						this.notifyUsers(c)
					}
				}
				if cmd.CmdType == game.CMDTYPE_DISCARD {
					result, _ := this.Game.Discard(idx, cmd.CardList)
					if result == game.DISCARD_SUCCESS {
						fmt.Printf("discard succeeded. next turn: %v\n", this.Game.Turn)
						this.notifyUsers(c)
					} else {
						fmt.Printf("Wrong Operation, code: %v\n", result)
					}
				}
			case <- this.Game.End:
				fmt.Println("End")
				winners := this.Game.GetWinner()
				for i := range winners {
					winners[i] = this.Users[winners[i]].Id
				}
				cmd := game.CardCommand{ CmdType: game.CMDTYPE_WIN, WinnerList: winners }
				this.notifyUsers(user.Command{ Id: 0, Command: cmd.ToMessage() })
				break PollLoop
		}
	}
}