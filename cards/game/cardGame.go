package game

import (
	"blackA/cards"
	"blackA/cards/pattern"
)

type ICardGame interface {
	PlayerNumber() int
	Compare(a, b *pattern.CardPattern) (int, bool)
	Init()
	DealCards() [][]cards.Card
}