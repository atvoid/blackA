package rule

import (
	"blackA/cards"
	"blackA/cards/pattern"
)

type ICardRule interface {
	PlayerNumber() int
	Compare(a, b *pattern.CardPattern) (int, bool)
	Init()
	DealCards() [][]cards.Card
}