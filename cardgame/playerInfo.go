package cardgame

import (
	"blackA/cards"
	"encoding/json"
)

type PlayerInfo struct {
	Group     int
	UserId    int
	Cards     []cards.Card
	DropCards []cards.Card
	IsClear   bool // indicate next round dicarding, clear all previous operation
	OnTurn    bool
	IsWinner  bool
	Score     int
}

func (this *PlayerInfo) ToMessage() string {
	s, _ := json.Marshal(*this)
	return string(s)
}
