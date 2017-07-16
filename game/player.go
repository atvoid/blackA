package game

import (
	"blackA/cards"
)

type Player struct {
	Group	int
	Cards	[]cards.Card
}

func (this *Player) Clear() {
	this.Cards = nil
}