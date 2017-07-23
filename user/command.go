package user

import (
	"encoding/json"
	"fmt"
)

const (
	CMDTYPE_LOGIN = 0
	CMDTYPE_GAME = 1
	CMDTYPE_DISCONNECT = 2
	CMDTYPE_JOINROOM = 3
	CMDTYPE_PING = 4

	CMDRESULT_ROOMFULL = 100
)

type Command struct {
	Id			int
	UserId		int
	CmdType		int
	Command		string
}

func (this *Command) ToMessage() string{
	s, _ := json.Marshal(*this)
	return string(s)
}

func MakeJoinRoomResult(uid int, rid int) Command {
	cmd := Command{ UserId: uid, CmdType: CMDTYPE_JOINROOM }
	cmd.Command = fmt.Sprintf("%v", rid)
	return cmd
}
