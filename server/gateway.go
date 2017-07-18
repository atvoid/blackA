package server

import (
	"blackA/user"
	"blackA/room"
	"net"
)

type ServerRouter struct {
	CmdFromUser			chan user.Command
	CmdFromServer		chan user.Command
	userSession			UserSession
	roomSession			RoomSession
	endSig				chan bool
}

var GlobalRouter ServerRouter = ServerRouter{
	CmdFromServer: make(chan user.Command, 100),
	CmdFromUser: make(chan user.Command, 100),
	userSession: GlobalUserSession,
	roomSession: GlobalRoomSession,
}

func (this *ServerRouter) AddUser(userid int, name string, conn *net.Conn) {
	ok := this.userSession.Login(userid, name, conn)
	if (ok == USERSESSION_LOGIN_SUCCESS) {
		u := this.userSession[userid]
		u.UserInput = this.CmdFromUser
	}
}

func (this *ServerRouter) AddRoom() {
	r := room.CreateRoom(this.CmdFromServer)
	this.roomSession[r.Id] = &r
}

func (this *ServerRouter) StartRouting() {
	RoutingLoop:
	for {
		select {
			case c := <- this.CmdFromUser:
				if c.CmdType == user.CMDTYPE_GAME {
					this.roomSession[this.userSession[c.UserId].RoomId].MsgReceiver <- c
				}
			case c := <- this.CmdFromServer:
				this.userSession[c.Id].ServerInput <- c
			case <- this.endSig:
				break RoutingLoop
		}
	}
}