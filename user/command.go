package user

import (
	"encoding/json"
)

const (
	CMDTYPE_LOGIN = 0
	CMDTYPE_GAME = 1
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