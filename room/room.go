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
	userCount	int
}

var uid int = 0

func CreateRoom(msgSender chan user.Command) Room {
	uid++
	room := Room{ Id: uid, MsgReceiver: make(chan user.Command, 20)}
	room.Game = game.CreateCardGame(0)
	room.Users = make([]int, len(room.Game.Players))
	room.MsgSender = msgSender
	room.userCount = 0;
	return room
}

func (this *Room) IsFull() bool {
	fmt.Printf("Room %v, Users: %v, Cap: %v\n", this.Id, this.userCount, len(this.Game.Players))
	if this.userCount < len(this.Game.Players) {
		return false
	}
	return true
}

func (this *Room) AddUser(id int) bool {
	if (!this.IsFull()) {
		//this.Users = append(this.Users, id)//append(this.Users, user.CreateUser(id, name, conn))
		this.userCount++
		for i, v := range this.Users {
			if v == 0 {
				this.Users[i] = id
				break
			}
		}
		//u := &this.Users[len(this.Users)-1]
		//u.UserInput = this.msgReceiver
		//go u.ReceiveMsg()
		//go u.HandleConnection()
		this.notifyStatus()
		if this.IsFull() {
			go this.Start()
		}
		return true
	} else {
		return false
	}
}

func (this *Room) RemoveUser(id int) {
	for i, v := range this.Users {
		if v == id {
			this.userCount--
			fmt.Printf("removed %v.\n", v)
			this.Users[i] = 0
			break
		}
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
	for _, id := range this.Users {
		if user.IsValidId(id) {
			cmd.UserId = id
			this.MsgSender <- cmd
			//fmt.Println(cmd.ToMessage())
		}
	}
}

func (this *Room) notifyStatus() {
	for i, id := range this.Users {
		if user.IsValidId(id) {
			msg := this.Game.GetStatus(i)
			msg.UserId = id
			for idx := range msg.PlayerList {
				msg.PlayerList[idx].UserId = this.Users[idx]
			}
			this.MsgSender <- user.Command{ UserId: id, CmdType: user.CMDTYPE_GAME, Command: msg.ToMessage() }
			// fmt.Println(msg.ToMessage())
		}
	}
}

func (this *Room) Start() {
	fmt.Printf("room %v started.\n", this.Id)
	// clear message
	ll := len(this.MsgReceiver)
	for i := 0; i < ll; i++ {
		<- this.MsgReceiver
	}
	this.Game.Start()
	this.notifyStatus()
	PollLoop:
	for {
		select {
			case c := <- this.MsgReceiver:
				if c.CmdType == user.CMDTYPE_DISCONNECT {
					// this.RemoveUser(c.UserId)
					break PollLoop
				}
				var cmd game.CardCommand
				idx := this.getUserIndex(c.UserId)
				json.Unmarshal([]byte(c.Command), &cmd)

				if cmd.CmdType == game.CMDTYPE_PASS {
					result := this.Game.Pass(idx)
					fmt.Printf("user %v pass.\n", c.UserId)
					if (result) {
						cmd.UserId = c.UserId
						this.notifyUsers(user.Command{ CmdType: user.CMDTYPE_GAME, Command: cmd.ToMessage()})
					}
					this.notifyStatus()
				} else if cmd.CmdType == game.CMDTYPE_DISCARD {
					fmt.Printf("Turn:%v, user %v, index %v discarding.\n", this.Game.Turn, c.UserId, idx)
					result, _ := this.Game.Discard(idx, cmd.CardList)
					if result == game.DISCARD_SUCCESS {
						fmt.Printf("discard succeeded. next turn: %v\n", this.Game.Turn)
						cmd.UserId = c.UserId
						this.notifyUsers(user.Command{ CmdType: user.CMDTYPE_GAME, Command: cmd.ToMessage()})
					} else {
						fmt.Printf("Wrong Operation, code: %v\n", result)
					}
					this.notifyStatus()
				} else if cmd.CmdType == game.CMDTYPE_INFO {
					this.notifyStatus()
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
	fmt.Printf("room %v ended.\n", this.Id)
	this.Clear()
}

func (this *Room) Clear() {
	this.Game.Clear()
}