package server

import (
	"blackA/user"
	"net"
)

const (
	USERSESSION_LOGIN_SUCCESS     = 0
	USERSESSION_RECONNECT_SUCCESS = 1
)

type UserSession map[int]*user.User

var GlobalUserSession UserSession = UserSession{}

func (this UserSession) Login(userid int, name string, conn *net.Conn) int {
	_, ok := this[userid]
	if ok {
		if this[userid].Disconnected {
			this[userid].Disconnected = false

			return USERSESSION_RECONNECT_SUCCESS
		}
		return -1
	} else {
		this[userid] = user.CreateUser(userid, name, conn)
		return USERSESSION_LOGIN_SUCCESS
	}
}
