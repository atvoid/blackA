package pattern

import (
	"blackA/cards"
	"sort"
)

const (
	BLACKAPATTERN_SINGLE = 1001
	BLACKAPATTERN_DOUBLE = 1002
	BLACKAPATTERN_BOMB = 1003
	BLACKAPATTERN_NUKE = 1004
	BLACKAPATTERN_STRAIGHT = 1005
	BLACKAPATTERN_DOUBLESTRAIGHT = 1006
	BLACKAPATTERN_DOUBLEKING = 1007
	BLACKAPATTERN_DOUBLEBLACKA = 1008
)

type blackAList []cards.Card;

func (this blackAList) Len() int {
	return len(this)
}

func (this blackAList) Less(i, j int) bool {
	return this[i].CardNumber < this[j].CardNumber
}

func (this blackAList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func blackAPattern_Single(list []cards.Card) int {
	if len(list) == 1 {
		return BLACKAPATTERN_SINGLE
	}
	return PATTERN_INVALID
}

func blackAPattern_Double(list []cards.Card) int {
	if len(list) == 2 {
		if (list[0].CardNumber == 1 &&
			list[1].CardNumber == 1 &&
			((list[0].CardType == cards.CARDTYPE_SPADE && list[1].CardType == cards.CARDTYPE_CLUB) ||
			(list[1].CardType == cards.CARDTYPE_SPADE && list[0].CardType == cards.CARDTYPE_CLUB))) {
			return PATTERN_INVALID
		}
		if (list[0].IsNormalCard() && list[1].IsNormalCard() && list[0].CardNumber == list[1].CardNumber) {
			return BLACKAPATTERN_DOUBLE
		}
	}
	return PATTERN_INVALID
}

func blackAPattern_Bomb(list []cards.Card) int {
	if len(list) == 3 {
		if (
			list[0].IsNormalCard() &&
			list[1].IsNormalCard() &&
			list[2].IsNormalCard() &&
			list[0].CardNumber == list[1].CardNumber &&
			list[1].CardNumber == list[2].CardNumber) {
			return BLACKAPATTERN_BOMB
		}
	}
	return PATTERN_INVALID
}

func blackAPattern_Nuke(list []cards.Card) int {
	if len(list) == 4 {
		if (
			list[0].IsNormalCard() &&
			list[1].IsNormalCard() &&
			list[2].IsNormalCard() &&
			list[3].IsNormalCard() &&
			list[0].CardNumber == list[1].CardNumber &&
			list[1].CardNumber == list[2].CardNumber &&
			list[2].CardNumber == list[3].CardNumber) {
			return BLACKAPATTERN_NUKE
		}
	}
	return PATTERN_INVALID
}

func blackAPattern_Straight(list []cards.Card) int {
	if len(list) >= 3 {
		sort.Sort(blackAList(list))
		if (list[0].CardNumber == 1) {
			if (list[1].CardNumber == 2) {
				// pattern A,2,3,4,5...
				for i, l := 2, len(list); i < l; i++ {
					if list[i].CardNumber - list[i-1].CardNumber != 1 {
						return PATTERN_INVALID
					}
				}
				return BLACKAPATTERN_STRAIGHT
			} else if (list[len(list)-1].CardNumber == 13) {
				// pattern ...,J,Q,K,A
				for i, l := 2, len(list); i < l; i++ {
					if list[i].CardNumber - list[i-1].CardNumber != 1 {
						return PATTERN_INVALID
					}
				}
				return BLACKAPATTERN_STRAIGHT
			} else {
				return PATTERN_INVALID
			}
		} else if list[0].CardNumber == 0 {
			return PATTERN_INVALID
		} else {
			// pattern ...6,7,8,9,10...
			for i, l := 1, len(list); i < l; i++ {
				if list[i].CardNumber - list[i-1].CardNumber != 1 {
					return PATTERN_INVALID
				}
			}
			return BLACKAPATTERN_STRAIGHT
		}
	}
	return PATTERN_INVALID
}

func blackAPattern_DoubleStraight(list []cards.Card) int {
	ll := len(list)
	if ll >= 6 && ll % 2 == 0 {
		sort.Sort(blackAList(list))
		if (list[1].CardNumber == 1) {
			if (list[3].CardNumber == 2) {
				// pattern AA,22,33,44,55...
				if (list[0].CardNumber != 1 || list[2].CardNumber != 2) {
					return PATTERN_INVALID
				}
				for i, l := 5, ll; i < l; i += 2 {
					if list[i].CardNumber - list[i-2].CardNumber != 1 || list[i].CardNumber != list[i-1].CardNumber {
						return PATTERN_INVALID
					}
				}
				return BLACKAPATTERN_DOUBLESTRAIGHT
			} else if (list[ll-1].CardNumber == 13) {
				// pattern ...,JJ,QQ,KK,AA
				if (list[0].CardNumber != 1 || list[ll-2].CardNumber != 13) {
					return PATTERN_INVALID
				}
				for i, l := 5, ll; i < l; i += 2 {
					if list[i].CardNumber - list[i-2].CardNumber != 1 || list[i].CardNumber != list[i-1].CardNumber {
						return PATTERN_INVALID
					}
				}
				return BLACKAPATTERN_DOUBLESTRAIGHT
			} else {
				return PATTERN_INVALID
			}
		} else if list[0].CardNumber == 0 {
			return PATTERN_INVALID
		} else {
			// pattern ...6,7,8,9,10...
			if (list[0].CardNumber != list[1].CardNumber) {
				return PATTERN_INVALID
			}
			for i, l := 3, ll; i < l; i += 2 {
				if list[i].CardNumber - list[i-2].CardNumber != 1 || list[i].CardNumber != list[i-1].CardNumber{
					return PATTERN_INVALID
				}
			}
			return BLACKAPATTERN_DOUBLESTRAIGHT
		}
	}
	return PATTERN_INVALID
}

func blackAPattern_DoubleKing(list []cards.Card) int {
	if len(list) == 2 {
		if (
			!list[0].IsNormalCard() &&
			!list[1].IsNormalCard()) {
			return BLACKAPATTERN_DOUBLEKING
		}
	}
	return PATTERN_INVALID
}

func blackAPattern_DoubleBlackA(list []cards.Card) int {
	if len(list) == 2 {
		if (
			list[0].CardNumber == 1 &&
			list[1].CardNumber == 1 &&
			((list[0].CardType == cards.CARDTYPE_SPADE && list[1].CardType == cards.CARDTYPE_CLUB) ||
			(list[1].CardType == cards.CARDTYPE_SPADE && list[0].CardType == cards.CARDTYPE_CLUB))) {
			return BLACKAPATTERN_DOUBLEBLACKA
		}
	}
	return PATTERN_INVALID
}

func GenerateBlackAPatternMatcher() func([]cards.Card) CardPattern {
	return CreatePatternClassifier([]func([]cards.Card) int{
		blackAPattern_Bomb,
		blackAPattern_Double,
		blackAPattern_DoubleBlackA,
		blackAPattern_DoubleKing,
		blackAPattern_Single,
		blackAPattern_Nuke,
		blackAPattern_Straight,
		blackAPattern_DoubleStraight})
}
