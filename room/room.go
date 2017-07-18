package room

import (
	"blackA/user"
	"blackA/game"
	"encoding/json"
	"fmt"
)

type Room struct {
	Id			int
	Game		game.CardGame
	Users		[]int
	MsgSender	chan user.Command
	MsgReceiver chan user.Command
}

var uid int = 0

func CreateRoom(msgSender chan user.Command) Room {
	uid++
	room := Room{ Id: uid, MsgReceiver: make(chan user.Command, 20)}
	room.Game = game.CreateCardGame(0)
	room.Users = make([]int, 0, len(room.Game.Players))
	room.MsgSender = msgSender
	return room
}

func (this *Room) IsFull() bool {
	fmt.Printf("Users: %v, Cap: %v\n", len(this.Users), cap(this.Users))
	if len(this.Users) < cap(this.Users) {
		return false
	}
	return true
}

func (this *Room) AddUser(id int) bool {
	if (!this.IsFull()) {
		this.Users = append(this.Users, id)//append(this.Users, user.CreateUser(id, name, conn))
		//u := &this.Users[len(this.Users)-1]
		//u.UserInput = this.msgReceiver
		//go u.ReceiveMsg()
		//go u.HandleConnection()
		return true
	} else {
		return false
	}
}

func (this *Room) getUserIndex(id int) int {
	for i, v := range this.Users {
		if v == id {
			return i
		}
	}
	panic("No Such User")
}

func (this *Room) notifyUsers(cmd user.Command) {
	for id := range this.Users {
		cmd.UserId = id
		this.MsgSender <- cmd
		//fmt.Println(cmd.ToMessage())
	}
}

func (this *Room) notifyStatus() {
	for i, id := range this.Users {
		msg := this.Game.GetStatus(i)
		msg.UserId = id
		for idx, p := range msg.PlayerList {
			p.UserId = this.Users[idx]
		}
		this.MsgSender <- user.Command{ Id: 0, Command: msg.ToMessage() }
		// fmt.Println(msg.ToMessage())
	}
}

func (this *Room) Start() {
	this.Game.Start()
	this.notifyStatus()
	PollLoop:
	for {
		select {
			case c:= <- this.MsgReceiver:
				fmt.Println(c.ToMessage())
				var cmd game.CardCommand
				idx := this.getUserIndex(c.UserId)
				json.Unmarshal([]byte(c.Command), &cmd)
				if cmd.CmdType == game.CMDTYPE_PASS {
					result := this.Game.Pass(idx)
					if (result) {
						cmd.UserId = c.UserId
						this.notifyUsers(user.Command{ Command: cmd.ToMessage()})
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
					winners[i] = this.Users[winners[i]]
				}
				cmd := game.CardCommand{ CmdType: game.CMDTYPE_WIN, WinnerList: winners }
				this.notifyUsers(user.Command{ Command: cmd.ToMessage() })
				break PollLoop
		}
	}
}