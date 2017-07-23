package server

import (
	"blackA/user"
	"blackA/room"
	"net"
	"fmt"
	"strconv"
	"blackA/logging"
)

var area string = "Router"

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
		go u.HandleConnection()
	}
}

func (this *ServerRouter) AddRoom() int {
	r := room.CreateRoom(this.CmdFromServer)
	this.roomSession[r.Id] = &r
	return r.Id
}

func (this *ServerRouter) JoinRoom(rid int, uid int) bool {
	result := this.roomSession[rid].AddUser(uid)
	if result {
		this.userSession[uid].RoomId = rid;
	}
	return result;
}

func (this *ServerRouter) StartRouting() {
	RoutingLoop:
	for {
		select {
			case c := <- this.CmdFromUser:
				_, has := this.userSession[c.UserId]
				if !has {
					logging.LogInfo(area, fmt.Sprintf("user %v does not exist\n", c.UserId))
				}
				if c.CmdType == user.CMDTYPE_PING {
					// do nothing
				} else if c.CmdType == user.CMDTYPE_GAME {
					this.roomSession[this.userSession[c.UserId].RoomId].MsgReceiver <- c
				} else if c.CmdType == user.CMDTYPE_DISCONNECT {
					logging.LogInfo(area, fmt.Sprintf("user %v disconnected\n", c.UserId))
					room, ok := this.roomSession[this.userSession[c.UserId].RoomId]
					if ok {
						room.MsgReceiver <- c
						room.RemoveUser(c.UserId)
					}
					delete(this.userSession, c.UserId)
				} else if c.CmdType == user.CMDTYPE_JOINROOM {
					num, err := strconv.ParseInt(c.Command, 10, 32)
					if err != nil {
						continue
					}
					rid := int(num)
					// room id is 0, join a random empty room
					if rid == 0 {
						succ := false;
						for i := range this.roomSession {
							if this.JoinRoom(i, c.UserId) {
								succ = true;
								this.sendJoinRoomResult(c.UserId, i)
								break
							}
						}
						// no valid room, join a new room
						if !succ {
							rid = this.AddRoom()
							logging.LogInfo(area, fmt.Sprintf("created room %v\n", rid))
							this.JoinRoom(rid, c.UserId)
							this.sendJoinRoomResult(c.UserId, rid)
						}
					} else {
						_, ok := this.roomSession[rid]
						if ok {
							if this.JoinRoom(rid, c.UserId) {
								// this.userSession[c.UserId].ServerInput <- user.Command{ UserId: c.UserId, CmdType: user.CMDRESULT_ROOMFULL }
								this.sendJoinRoomResult(c.UserId, rid)
								break;
							}
						}
						this.sendJoinRoomResult(c.UserId , 0)
					}
				}
			case c := <- this.CmdFromServer:
				logging.LogInfo(area, fmt.Sprintf("server command to %v\n", c.UserId))
				this.userSession[c.UserId].ServerInput <- c
			case <- this.endSig:
				break RoutingLoop
		}
	}
}

func (this *ServerRouter) sendJoinRoomResult(uid, rid int) {
	cmd := user.MakeJoinRoomResult(uid, rid)
	this.userSession[uid].ServerInput <- cmd
}