package server

import (
	"blackA/logging"
	"blackA/room"
	"blackA/user"
	"fmt"
	"net"
)

var area string = "Router"

type ServerRouter struct {
	CmdFromUser   chan user.Command
	CmdFromServer chan user.Command
	userSession   UserSession
	roomSession   RoomSession
	endSig        chan bool
	area          string
}

var GlobalRouter ServerRouter = ServerRouter{
	CmdFromServer: make(chan user.Command, 100),
	CmdFromUser:   make(chan user.Command, 100),
	userSession:   GlobalUserSession,
	roomSession:   GlobalRoomSession,
}

func (this *ServerRouter) AddUser(userid int, name string, conn *net.Conn) {
	ok := this.userSession.Login(userid, name, conn)
	if ok == USERSESSION_LOGIN_SUCCESS {
		u := this.userSession[userid]
		u.UserInput = this.CmdFromUser
		go u.HandleConnection()
	} else if ok == USERSESSION_RECONNECT_SUCCESS {
		u := this.userSession[userid]
		u.ResetConnection(conn)
		u.UserInput = this.CmdFromUser
		go u.HandleConnection()
		logging.LogInfo_Detail(area, fmt.Sprintf("user %v reconnected.", userid))

		if u.RoomId > 0 {
			r, ok := this.roomSession[u.RoomId]
			if ok {
				r.MsgReceiver <- room.MakeUserCommandForRoom(userid, room.MakeRoomRequest_Reconnect(r.Id))
			}
		}
	}
}

func (this *ServerRouter) AddRoom() int {
	r := room.CreateRoom(this.CmdFromServer, nil) // TODO
	this.roomSession[r.Id] = &r
	go r.Start()
	return r.Id
}

func (this *ServerRouter) JoinRoom(rid int, uid int) bool {
	result := this.roomSession[rid].AddUser(uid)
	if result {
		this.userSession[uid].RoomId = rid
	}
	return result
}

func (this *ServerRouter) handlerUserCommand(c user.Command) {
	u, has := this.userSession[c.UserId]
	if !has {
		logging.LogInfo(area, fmt.Sprintf("user %v does not exist\n", c.UserId))
		return
	}

	switch c.CmdType {
	case user.CMDTYPE_GAME:
		r, ok := this.roomSession[u.RoomId]
		if ok {
			r.MsgReceiver <- c
		} else {
			logging.LogInfo_Detail(area, fmt.Sprintf("game cmd to room %v, not exist. user: %v.", u.RoomId, u.Id))
		}
	case user.CMDTYPE_DISCONNECT:
		logging.LogInfo(area, fmt.Sprintf("user %v disconnected\n", c.UserId))
		r, ok := this.roomSession[this.userSession[c.UserId].RoomId]
		if ok {
			this.userSession[c.UserId].Disconnected = true
			r.MsgReceiver <- room.MakeUserCommandForRoom(c.UserId, room.MakeRoomRequest_Disconnect(r.Id))
		} else {
			delete(this.userSession, c.UserId)
		}
	case user.CMDTYPE_ROOM:
		if u.RoomId == 0 && c.RoomId == 0 {
			// join a random room
			for _, v := range this.roomSession {
				if v != nil && !v.IsFull() {
					v.MsgReceiver <- c
				}
			}
			rid := this.AddRoom()
			this.roomSession[rid].MsgReceiver <- c
		} else if u.RoomId == 0 && c.RoomId > 0 {
			// join a specific room
			r, ok := this.roomSession[c.RoomId]
			if ok {
				r.MsgReceiver <- c
			} else {
				rid := this.AddRoom()
				this.roomSession[rid].MsgReceiver <- c
			}
		} else {
			// send message to the room user exists
			r, ok := this.roomSession[u.RoomId]
			if ok {
				r.MsgReceiver <- c
			} else {
				logging.LogError(area, fmt.Sprintf("room %v doesn't exist. user: %v", u.RoomId, c.UserId))
			}
		}
	}
}

func (this *ServerRouter) StartRouting() {
RoutingLoop:
	for {
		select {
		case c := <-this.CmdFromUser:
			if c.CmdType == user.CMDTYPE_INTERNAL_ROOMEMPTY {
				this.roomSession[c.RoomId].Dispose()
				delete(this.roomSession, c.RoomId)
			}
			this.handlerUserCommand(c)
		case c := <-this.CmdFromServer:
			logging.LogInfo_Detail(this.area, fmt.Sprintf("server command: %v\n", c.ToMessage()))
			if c.CmdType == user.CMDTYPE_INTERNAL_ROOMEMPTY {
				r, ok := this.roomSession[c.RoomId]
				if ok && r != nil {
					r.Dispose()
					delete(this.roomSession, c.RoomId)
				}
			} else {
				this.userSession[c.UserId].ServerInput <- c
			}
		case <-this.endSig:
			break RoutingLoop
		}
	}
}
