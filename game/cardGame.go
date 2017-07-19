package game

import (
	"blackA/cards/rule"
	"blackA/cards"
	"blackA/cards/pattern"
	"sort"
	"fmt"
)

const (
	DISCARD_INVALIDPATTERN = 0
	DISCARD_SUCCESS = 1
	DISCARD_SMALLER = 2
	DISCARD_WRONGTURN = 3
	DISCARD_NOTEXIST = 4

	BLACKA_GROUP = 0
	NONBLACKA_GROUP = 1
)

type dropCards struct {
	PlayerId	int
	Cards		pattern.CardPattern
}

type CardGame struct {
	rule		rule.ICardRule
	DropList	[]dropCards
	Players		[]Player
	Turn		int
	End			chan bool
	wind		bool
	winGroup	int
}

func CreateCardGame(startTurn int) CardGame {
	game := CardGame{
		rule: rule.ICardRule(&rule.BlackAGame{}),
		Turn: startTurn,
	}
	game.Players = make([]Player, game.rule.PlayerNumber())
	game.DropList = make([]dropCards, 0, 60)
	game.End = make(chan bool)
	game.wind = false;
	return game
}

func (this *CardGame) Start() {
	this.rule.Init()
	deal := this.rule.DealCards()
	for i, v := range deal {
		this.Players[i].Cards = v
		if (this.hasCards(i, []cards.Card{{ CardNumber:1, CardType: cards.CARDTYPE_CLUB }}) ||
			this.hasCards(i, []cards.Card{{ CardNumber:1, CardType: cards.CARDTYPE_SPADE }})) {
				this.Players[i].Group = BLACKA_GROUP
		} else {
			this.Players[i].Group = NONBLACKA_GROUP
		}
	}
}

func (this *CardGame) Clear() {
	for i := range this.Players {
		this.Players[i].Cards = make([]cards.Card, 0)
		this.Players[i].Group = 0
	}
	this.DropList = make([]dropCards, 0)
	this.wind = false
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

func (this *CardGame) nextTurn() {
	total := this.rule.PlayerNumber()
	this.Turn = (this.Turn + 1) % total
	for len(this.Players[this.Turn].Cards) == 0 {
		// if there is no one discarding and last turn player has no card, then next one take the turn.
		if this.DropList[len(this.DropList)-1].PlayerId == this.Turn {
			this.wind = true
		}
		this.Turn = (this.Turn + 1) % total
	}
}

func (this *CardGame) Pass(playerNumber int) bool {
	if (this.Turn != playerNumber) {
		return false;
	}
	ll := len(this.DropList)
	if (ll > 0 && this.DropList[ll-1].PlayerId == playerNumber) {
		return false;
	}
	this.nextTurn()
	return true;
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

	blackAWin := true
	nonAWin := true
	for _, v := range this.Players {
		if len(v.Cards) != 0 {
			if v.Group == BLACKA_GROUP {
				blackAWin = false;
			} else {
				nonAWin = false;
			}
		}
	}
	if blackAWin || nonAWin {
		if blackAWin {
			fmt.Printf("black win.\n")
			this.winGroup = BLACKA_GROUP
		}
		if nonAWin {
			fmt.Printf("non-black win.\n")
			this.winGroup = NONBLACKA_GROUP
		}
		this.End <- true
	} else {
		this.nextTurn()
	}
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
	if (ll == 0 || this.wind || this.DropList[ll-1].PlayerId == playerNumber) {
		this.wind = false
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

func (this *CardGame) GetWinner() []int {
	winner := make([]int, 0, this.rule.PlayerNumber())
	for i, v := range this.Players {
		if v.Group == this.winGroup {
			winner = append(winner, i)
		}
	}
	return winner
}

func (this *CardGame) GetStatus(pIdx int) CardCommand {
	ans := make([]Player, len(this.Players))
	for i, v := range this.Players {
		if i != pIdx {
			ans[i].Cards = make([]cards.Card, len(v.Cards))
		} else {
			ans[i].Cards = v.Cards
		}
		ans[i].OnTurn = this.Turn == i
	}
	return CardCommand{ CmdType: CMDTYPE_INFO, PlayerList: ans }
}