package pattern

import (
	"blackA/cards"
)

type CardPattern struct {
	PatternType	int
	CardList	[]cards.Card
}

const (
	PATTERN_INVALID = 0
)

func CreatePatternClassifier(list []func([]cards.Card) int) func([]cards.Card) CardPattern {
	return func(cardList []cards.Card) CardPattern {
		for _, f := range list {
			pattern := f(cardList)
			if (pattern != PATTERN_INVALID) {
				return CardPattern{ PatternType: pattern, CardList: cardList }
			}
		}
		return CardPattern{ PatternType: PATTERN_INVALID, CardList: cardList }
	}
}