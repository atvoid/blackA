package game

import (
	"blackA/cards"
)

type Player struct {
	Cards	[]cards.Card
}

func (this *Player) Clear() {
	this.Cards = nil
}