package user

import (
	"encoding/json"
)

const (
	CMDTYPE_LOGIN      = 0
	CMDTYPE_GAME       = 1
	CMDTYPE_DISCONNECT = 2
	CMDTYPE_PING       = 4
	CMDTYPE_ROOM       = 5

	// internal Command
	CMDTYPE_INTERNAL_ROOMEMPTY = 1000
	CMDTYPE_INTERNAL_LEAVEROOM = 1001
	CMDTYPE_INTERNAL_JOINROOM  = 1002
)

type Command struct {
	Id      int
	UserId  int
	RoomId  int
	CmdType int
	Command string
}

func (this *Command) ToMessage() string {
	s, _ := json.Marshal(*this)
	return string(s)
}

/*
func MakeJoinRoomResult(uid int, rid int) Command {
	cmd := Command{UserId: uid, CmdType: CMDTYPE_JOINROOM}
	cmd.Command = fmt.Sprintf("%v", rid)
	return cmd
}
*/
