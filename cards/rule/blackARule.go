package rule

import (
	"blackA/cards/pattern"
	"blackA/cards"
)

type BlackAGame struct {
	suite	cards.BlackASuite
	matcher	func([]cards.Card) pattern.CardPattern
}

func (this *BlackAGame) PlayerNumber() int {
	return 4
}

func (this *BlackAGame) Init() {
	this.suite = cards.CreateCardSuiteForBlackA()
	this.suite.Shuffle()
}

func (this *BlackAGame) DealCards() [][]cards.Card {
	num := this.PlayerNumber()
	ans := make([][]cards.Card, num)
	for i := range ans {
		ans[i] = make([]cards.Card, 0, len(this.suite.CardList))
	}
	for i, v := range this.suite.CardList {
		ans[i%num] = append(ans[i%num], v)
	}
	return ans
}

func (this *BlackAGame) GetPatternMatcher() func([]cards.Card) pattern.CardPattern {
	if (this.matcher == nil) {
		this.matcher = pattern.GenerateBlackAPatternMatcher()
	}
	return this.matcher
}

func numberCompare(a, b *cards.Card) int {
	getWeight := func(c *cards.Card) int {
		if c.CardType == cards.CARDTYPE_SPADE && c.CardNumber == 1 {
			return 99
		} else if c.CardType == cards.CARDTYPE_CLUB && c.CardNumber == 1 {
			return 98
		} else if c.CardType == cards.CARDTYPE_JOKER_L {
			return 97
		} else if c.CardType == cards.CARDTYPE_JOKER_S {
			return 96
		} else if c.CardNumber == 2 {
			return 95
		} else if c.CardNumber == 1 {
			return 94
		} else {
			return c.CardNumber
		}
	}
	wa, wb := getWeight(a), getWeight(b)
	return wa-wb;
}

func doubleCompare(a, b *cards.Card) int {
	getWeight := func(c *cards.Card) int {
		if c.CardNumber == 2 {
			return 95
		} else if c.CardNumber == 1 {
			return 94
		} else {
			return c.CardNumber
		}
	}
	wa, wb := getWeight(a), getWeight(b)
	return wa-wb;
}

// pattern is sorted
func straightCompare(a, b *pattern.CardPattern) int {
	ll := len(a.CardList)
	if (a.CardList[0].CardNumber == 1) {
		// a is A2345
		if (a.CardList[1].CardNumber == 2) {
			if (b.CardList[0].CardNumber == 1 && b.CardList[1].CardNumber == 2) {
				return 0
			} else {
				return -1
			}
		} else {
			// a is JQKA
			if (b.CardList[0].CardNumber == 1 && b.CardList[ll-1].CardNumber == 13) {
				return 0
			} else {
				return 1
			}
		}
	// a has no A
	} else {
		// b is A2345
		if (b.CardList[0].CardNumber == 1 && b.CardList[1].CardNumber == 2) {
			return 1
		// b is JQKA
		} else if (b.CardList[0].CardNumber == 1 && b.CardList[ll-1].CardNumber == 13) {
			return -1
		} else {
		// b has no A
			return numberCompare(&a.CardList[0], &b.CardList[0])
		}
	}
}

// pattern is sorted
func doubleStraightCompare(a, b *pattern.CardPattern) int {
	ll := len(a.CardList)
	if (a.CardList[0].CardNumber == 1) {
		// a is AA22334455
		if (a.CardList[2].CardNumber == 2) {
			if (b.CardList[0].CardNumber == 1 && b.CardList[2].CardNumber == 2) {
				return 0
			} else {
				return -1
			}
		} else {
			// a is JJQQKKAA
			if (b.CardList[0].CardNumber == 1 && b.CardList[ll-1].CardNumber == 13) {
				return 0
			} else {
				return 1
			}
		}
	// a has no A
	} else {
		// b is AA22334455
		if (b.CardList[0].CardNumber == 1 && b.CardList[2].CardNumber == 2) {
			return 1
		// b is JJQQKKAA
		} else if (b.CardList[0].CardNumber == 1 && b.CardList[ll-1].CardNumber == 13) {
			return -1
		} else {
		// b has no A
			return numberCompare(&a.CardList[0], &b.CardList[0])
		}
	}
}

func (this *BlackAGame) Compare(a, b *pattern.CardPattern) (int, bool) {
	if (b.PatternType == pattern.BLACKAPATTERN_SINGLE) {
		switch {
			case a.PatternType == pattern.BLACKAPATTERN_SINGLE:
				return numberCompare(&a.CardList[0], &b.CardList[0]), true
			case (a.PatternType == pattern.BLACKAPATTERN_DOUBLE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_STRAIGHT):
				return 0, false
			case (a.PatternType == pattern.BLACKAPATTERN_BOMB ||
				a.PatternType == pattern.BLACKAPATTERN_NUKE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEKING):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_DOUBLE) {
		switch {
			case a.PatternType == pattern.BLACKAPATTERN_DOUBLE:
				return doubleCompare(&a.CardList[0], &b.CardList[0]), true
			case (a.PatternType == pattern.BLACKAPATTERN_SINGLE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_STRAIGHT):
				return 0, false
			case (a.PatternType == pattern.BLACKAPATTERN_BOMB ||
				a.PatternType == pattern.BLACKAPATTERN_NUKE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEKING):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_STRAIGHT) {
		switch {
			case a.PatternType == pattern.BLACKAPATTERN_STRAIGHT && len(a.CardList) == len(b.CardList):
				return straightCompare(a, b), true
			case (a.PatternType == pattern.BLACKAPATTERN_SINGLE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLE):
				return 0, false
			case (a.PatternType == pattern.BLACKAPATTERN_BOMB ||
				a.PatternType == pattern.BLACKAPATTERN_NUKE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEKING):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT) {
		switch {
			case a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT && len(a.CardList) == len(b.CardList):
				return doubleStraightCompare(a, b), true
			case (a.PatternType == pattern.BLACKAPATTERN_SINGLE ||
				a.PatternType == pattern.BLACKAPATTERN_STRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLE):
				return 0, false
			case (a.PatternType == pattern.BLACKAPATTERN_BOMB ||
				a.PatternType == pattern.BLACKAPATTERN_NUKE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEKING):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_BOMB) {
		switch {
			case a.PatternType == pattern.BLACKAPATTERN_BOMB:
				return doubleCompare(&a.CardList[0], &b.CardList[0]), true
			case (a.PatternType == pattern.BLACKAPATTERN_SINGLE ||
				a.PatternType == pattern.BLACKAPATTERN_STRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLE):
				return -1, true
			case (a.PatternType == pattern.BLACKAPATTERN_NUKE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEKING):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_NUKE) {
		switch {
			case a.PatternType == pattern.BLACKAPATTERN_NUKE:
				return doubleCompare(&a.CardList[0], &b.CardList[0]), true
			case (a.PatternType == pattern.BLACKAPATTERN_SINGLE ||
				a.PatternType == pattern.BLACKAPATTERN_STRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_BOMB ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLE):
				return -1, true
			case (a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLEKING):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_DOUBLEKING) {
		switch {
			case (a.PatternType == pattern.BLACKAPATTERN_SINGLE ||
				a.PatternType == pattern.BLACKAPATTERN_STRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLESTRAIGHT ||
				a.PatternType == pattern.BLACKAPATTERN_BOMB ||
				a.PatternType == pattern.BLACKAPATTERN_NUKE ||
				a.PatternType == pattern.BLACKAPATTERN_DOUBLE):
				return -1, true
			case (a.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA):
				return 1, true
			default:
				return 0, false
		}
	} else if (b.PatternType == pattern.BLACKAPATTERN_DOUBLEBLACKA) {
		return -1, true
	}
	return 0, false
}