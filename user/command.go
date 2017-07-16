package user

import (
	"encoding/json"
)

type Command struct {
	Id			int
	UserId		int
	Command		string
}

func (this *Command) ToMessage() string{
	s, _ := json.Marshal(*this)
	return string(s)
}