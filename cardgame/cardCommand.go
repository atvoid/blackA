package cardgame

import (
	"blackA/cards"
	"encoding/json"
)

const (
	CMDTYPE_DISCARD = 0
	CMDTYPE_PASS    = 1
	CMDTYPE_INFO    = 2
	CMDTYPE_WIN     = 3
)

type CardCommand struct {
	CmdType    int
	UserId     int
	CardList   []cards.Card
	PlayerList []PlayerInfo
	WinnerList []int
}

func (this *CardCommand) ToMessage() string {
	s, _ := json.Marshal(*this)
	return string(s)
}
