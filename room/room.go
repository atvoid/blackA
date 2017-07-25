package room

import (
	"blackA/logging"
	"blackA/user"
	"encoding/json"
	"fmt"
)

var area string = "Room"

type Room struct {
	Id              int
	GameCreator     IGameCreator
	Game            IGame
	Users           []*UserInfo
	MsgSender       chan user.Command
	MsgReceiver     chan user.Command
	gameMsgReceiver chan GameCommand
	gameMsgSender   chan GameCommand
	disposeSig      chan bool
	userCount       int
	isStarted       bool
	nextGameData    interface{}
	logArea         string
}

var uid int = 10000

func CreateRoom(msgSender chan user.Command, gameCreator IGameCreator) Room {
	uid++
	room := Room{Id: uid, MsgReceiver: make(chan user.Command, 40)}
	room.isStarted = false
	room.GameCreator = gameCreator
	room.Users = make([]*UserInfo, gameCreator.GetPlayerCapacity())
	room.MsgSender = msgSender
	room.userCount = 0
	room.nextGameData = nil
	room.logArea = area + fmt.Sprintf("_%v", room.Id)

	room.gameMsgReceiver = make(chan GameCommand, 30)
	room.gameMsgSender = make(chan GameCommand, 30)
	return room
}

func (this *Room) resetAll() {
	for _, v := range this.Users {
		if v != nil {
			v.Ready = false
		}
	}
	this.nextGameData = nil
	this.isStarted = false
}

func (this *Room) IsFull() bool {
	logging.LogInfo_Detail(this.logArea, fmt.Sprintf("Room %v, Users: %v, Cap: %v\n", this.Id, this.userCount, this.GameCreator.GetPlayerCapacity()))
	if this.userCount < this.GameCreator.GetPlayerCapacity() {
		return false
	}
	return true
}

func (this *Room) AddUser(id int) bool {
	if !this.IsFull() && this.getUserIndex(id) == -1 {
		this.userCount++
		for i, v := range this.Users {
			if v == nil {
				u := NewUserInfo(id)
				this.Users[i] = &u
				logging.LogInfo_Normal(this.logArea, fmt.Sprintf("added user %v to room %v.\n", id, this.Id))
				break
			}
		}
		return true
	} else {
		return false
	}
}

func (this *Room) RemoveUser(id int) {
	for i, v := range this.Users {
		if v != nil && v.UserId == id {
			this.userCount--
			logging.LogInfo_Normal(this.logArea, fmt.Sprintf("removed user %v from room %v.\n", v.UserId, this.Id))
			this.Users[i] = nil
			break
		}
	}
}

func (this *Room) getUserIndex(id int) int {
	for i, v := range this.Users {
		if v != nil && v.UserId == id {
			return i
		}
	}
	return -1
}

func (this *Room) allReady() bool {
	for _, v := range this.Users {
		if v == nil {
			return false
		} else if !v.Ready {
			return false
		}
	}
	return true
}

func (this *Room) getRoomStatus() []*UserInfo {
	uInfo := make([]*UserInfo, len(this.Users))
	for i, v := range this.Users {
		if v != nil {
			u := NewUserInfo(v.UserId)
			uInfo[i] = &u
			uInfo[i].Connected = v.Connected
			uInfo[i].PlayerInfo = nil
			uInfo[i].Ready = v.Ready
		}
	}
	return uInfo
}

func (this *Room) transformGameCommand(gCmd GameCommand) *user.Command {
	if this.Users[gCmd.PlayerIndex] == nil {
		logging.LogError(area, fmt.Sprintf("room index %v has no user.", gCmd.PlayerIndex))
		return nil
	}
	var c user.Command
	switch gCmd.CmdType {
	case GAMECMD_NOTIFY:
		uInfo := this.getRoomStatus()

		for i, v := range gCmd.PlayerInfo {
			if uInfo[i] != nil {
				uInfo[i].PlayerInfo = &v
			}
		}
		c = MakeUserCommandForRoom(this.Users[gCmd.PlayerIndex].UserId, MakeRoomResponse_GameStatus(this.Id, uInfo, gCmd.GameInfo, this.isStarted))
	}
	return &c
}

func (this *Room) transformUserCommand(uCmd user.Command) *GameCommand {
	idx := this.getUserIndex(uCmd.UserId)
	if idx == -1 {
		logging.LogError(area, fmt.Sprintf("room has no userId %v.", uCmd.UserId))
		return nil
	}
	gCmd := MakeGameCommandRequest_Game(idx, uCmd.Command)
	return &gCmd
}

func (this *Room) notifyRoomStatus() {
	uInfos := this.getRoomStatus()
	for _, v := range this.Users {
		if v != nil {
			c := MakeUserCommandForRoom(v.UserId, MakeRoomResponse_RoomStatus(this.Id, uInfos, this.isStarted))
			this.MsgSender <- c
		}
	}
}

func (this *Room) startGame() {
	this.Game = this.GameCreator.CreateGame(this.nextGameData)
	this.Game.SetMsgSender(this.gameMsgReceiver)
	this.Game.SetMsgReceiver(this.gameMsgSender)
	// clear all ready flag
	for _, v := range this.Users {
		if v != nil {
			v.Ready = false
		}
	}
	// clear all game msg buffer
	ll := len(this.gameMsgSender)
	for i := 0; i < ll; i++ {
		<-this.gameMsgSender
	}
	this.Game.Start()
	this.isStarted = true
}

func (this *Room) endGame(nextGameData interface{}) {
	this.nextGameData = nextGameData
	this.Game = nil
	this.isStarted = false
}

func (this *Room) handleRoomCommand(uid int, rCmd RoomCommand) {
	switch rCmd.CmdType {
	case ROOMCMD_DISCONNECT:
		// if game does not start, will remove the user
		if !this.isStarted {
			this.RemoveUser(uid)
			this.notifyRoomStatus()
		} else {
			pIdx := this.getUserIndex(uid)
			if pIdx != -1 {
				this.Users[pIdx].Connected = false

				// notify game a user disconnected
				gCmd := MakeGameCommandRequest_Disconnect(pIdx)
				this.gameMsgSender <- gCmd
			}
		}
	case ROOMCMD_RECONNECT:
		// if game does not start, just notify all
		pIdx := this.getUserIndex(uid)
		if pIdx != -1 {
			if !this.isStarted {
				this.Users[pIdx].Connected = true
				this.notifyRoomStatus()
			} else {
				// notify game a user reconnected
				gCmd := MakeGameCommandRequest_Reconnect(pIdx)
				this.gameMsgSender <- gCmd
			}
		}
	case ROOMCMD_JOIN:
		if this.IsFull() {
			logging.LogInfo_Detail(this.logArea, "room is full.")
			this.MsgSender <- MakeUserCommandForRoom(uid, MakeRoomResponse_Join_Full(this.Id))
		} else if this.isStarted {
			logging.LogInfo_Detail(this.logArea, "room has started.")
			this.MsgSender <- MakeUserCommandForRoom(uid, MakeRoomResponse_Join_Started(this.Id))
		} else {
			// if this user is already in room, view it as success
			this.AddUser(uid)
			this.MsgSender <- MakeUserCommandForRoom(uid, MakeRoomResponse_Join_Success(this.Id))
			this.notifyRoomStatus()
		}
	case ROOMCMD_LEAVE:
		if this.isStarted {
			pIdx := this.getUserIndex(uid)
			if pIdx != -1 {
				logging.LogInfo_Detail(this.logArea, fmt.Sprintf("%v exited from a ongoing game.", uid))
				this.gameMsgSender <- MakeGameCommandRequest_Exit(pIdx)
				this.notifyRoomStatus()
			}
		} else {
			this.RemoveUser(uid)
			this.MsgSender <- MakeUserCommandForRoom(uid, MakeRoomResponse_Leave_Success(this.Id))
			// if the room is empty, notify upper layer to dispose this room
			if this.userCount == 0 {
				this.MsgSender <- user.Command{CmdType: user.CMDTYPE_INTERNAL_ROOMEMPTY, RoomId: this.Id}
			} else {
				this.notifyRoomStatus()
			}
		}
	case ROOMCMD_READY:
		if this.isStarted {
			logging.LogInfo_Detail(this.logArea, fmt.Sprintf("user:%v, invalid ready cmd, game has already started.", uid))
		} else {
			pIdx := this.getUserIndex(uid)
			if pIdx != -1 {
				this.Users[pIdx].Ready = true
				if this.allReady() {
					this.startGame()
				}
			}
		}
		this.notifyRoomStatus()
	case ROOMCMD_NOTREADY:
		if this.isStarted {
			logging.LogInfo_Detail(this.logArea, fmt.Sprintf("user:%v, invalid notready cmd, game has already started.", uid))
		} else {
			pIdx := this.getUserIndex(uid)
			if pIdx != -1 {
				this.Users[pIdx].Ready = false
			}
		}
		this.notifyRoomStatus()
	case ROOMCMD_ROOMSTATUS:
		uInfos := this.getRoomStatus()
		this.MsgSender <- MakeUserCommandForRoom(uid, MakeRoomResponse_RoomStatus(this.Id, uInfos, this.isStarted))
	case ROOMCMD_GAMESTATUS:
		if this.isStarted {
			pIdx := this.getUserIndex(uid)
			if pIdx != -1 {
				this.gameMsgSender <- MakeGameCommandRequest_Status(pIdx)
			}
		}
	}
}

func (this *Room) Dispose() {
	this.disposeSig <- true
}

func (this *Room) Start() {
	logging.LogInfo(area, fmt.Sprintf("room %v started.\n", this.Id))
	// clear message
	ll := len(this.MsgReceiver)
	for i := 0; i < ll; i++ {
		<-this.MsgReceiver
	}

PollLoop:
	for {
		select {
		case c := <-this.MsgReceiver:
			if c.CmdType == user.CMDTYPE_GAME {

				gCmd := this.transformUserCommand(c)
				if gCmd != nil {
					this.gameMsgSender <- *gCmd
				}

			} else if c.CmdType == user.CMDTYPE_ROOM {

				var rCmd RoomCommand
				err := json.Unmarshal([]byte(c.Command), &rCmd)
				if err != nil {
					logging.LogError(area, fmt.Sprintf("invalid user command:%v", c.ToMessage()))
				}
				this.handleRoomCommand(c.UserId, rCmd)

			}

		case gCmd := <-this.gameMsgReceiver:
			if gCmd.CmdType == GAMECMD_NOTIFY {
				c := this.transformGameCommand(gCmd)
				if c != nil {
					this.MsgSender <- *c
				}
			} else if gCmd.CmdType == GAMECMD_GAMEEND {
				this.endGame(gCmd.NextGameData)
			}

		case <-this.disposeSig:
			break PollLoop
		}
	}
	logging.LogInfo(area, fmt.Sprintf("room %v terminated.\n", this.Id))
}
