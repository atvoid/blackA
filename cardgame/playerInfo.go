package game

import (
	"blackA/cards"
	"encoding/json"
)

type PlayerInfo struct {
	Group  int
	UserId int
	Cards  []cards.Card
	OnTurn bool
}

func (this *PlayerInfo) Clear() {
	this.Cards = nil
}

func (this *PlayerInfo) ToMessage() string {
	s, _ := json.Marshal(this)
	return string(s)
}
