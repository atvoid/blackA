package game

import (
	"blackA/cards/rule"
	"blackA/cards"
	"blackA/cards/pattern"
	"sort"
)

const (
	DISCARD_INVALIDPATTERN = 0
	DISCARD_SUCCESS = 1
	DISCARD_SMALLER = 2
	DISCARD_WRONGTURN = 3
	DISCARD_NOTEXIST = 4
)

type dropCards struct {
	PlayerId	int
	Cards		pattern.CardPattern
}

type CardGame struct {
	rule		rule.ICardRule
	blackAGroup	[]int
	DropList	[]dropCards
	Players		[]Player
	Turn		int
}

func CreateCardGame() CardGame {
	game := CardGame{
		rule: rule.ICardRule(&rule.BlackAGame{}),
		Turn: 0,
	}
	game.Players = make([]Player, game.rule.PlayerNumber())
	game.blackAGroup = make([]int, 0, game.rule.PlayerNumber())
	game.DropList = make([]dropCards, 0, 60)
	return game
}

func (this *CardGame) Start() {
	this.rule.Init()
	deal := this.rule.DealCards()
	for i, v := range deal {
		this.Players[i].Cards = v
		if (this.hasCards(i, []cards.Card{{ CardNumber:1, CardType: cards.CARDTYPE_CLUB }}) ||
			this.hasCards(i, []cards.Card{{ CardNumber:1, CardType: cards.CARDTYPE_SPADE }})) {
				this.blackAGroup = append(this.blackAGroup, i)
		}
	}
}

func (this *CardGame) hasCards(playerNumer int, cardList []cards.Card) bool {
	sort.Sort(cards.CardList(this.Players[playerNumer].Cards))
	sort.Sort(cards.CardList(cardList))
	l1, l2 := len(this.Players[playerNumer].Cards), len(cardList)
	for i, j := 0, 0; i < l1 && j < l2; {
		for i < l1 && !this.Players[playerNumer].Cards[i].Equal(&cardList[j]) {
			i++
		}
		if (i == l1) {
			return false
		}
		i++
		j++
	}
	return true
}

func (this *CardGame) discard(playerNumber int, pat pattern.CardPattern) {
	for _, v := range pat.CardList {
		idx := 0
		for !this.Players[playerNumber].Cards[idx].Equal(&v) {
			idx++
		}
		this.Players[playerNumber].Cards = append(this.Players[playerNumber].Cards[:idx], this.Players[playerNumber].Cards[idx+1:]...)
	}
	this.DropList = append(this.DropList, dropCards{ PlayerId: playerNumber, Cards: pat })
}

func (this *CardGame) Discard(playerNumber int, cardList []cards.Card) (int, *pattern.CardPattern) {
	if (this.Turn != playerNumber) {
		return DISCARD_WRONGTURN, nil
	}
	if (!this.hasCards(playerNumber, cardList)) {
		return DISCARD_NOTEXIST, nil
	}
	pat := this.rule.GetPatternMatcher()(cardList)
	if (pat.PatternType == pattern.PATTERN_INVALID) {
		return DISCARD_INVALIDPATTERN, nil
	}
	ll := len(this.DropList)
	if (ll == 0 || this.DropList[ll-1].PlayerId == playerNumber) {
		this.discard(playerNumber, pat)
		return DISCARD_SUCCESS, &pat
	}
	val, ok := this.rule.Compare(&pat, &this.DropList[ll-1].Cards)
	if (ok && val > 0) {
		this.discard(playerNumber, pat)
		return DISCARD_SUCCESS, &pat
	}
	return DISCARD_SMALLER, nil
}