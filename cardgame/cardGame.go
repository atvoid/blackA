package cardgame

import (
	"blackA/cards"
	"blackA/cards/pattern"
	"blackA/cards/rule"
	"blackA/logging"
	"blackA/room"
	"encoding/json"
	"fmt"
	"sort"
)

var area string = "CardGame"

const (
	DISCARD_INVALIDPATTERN = 0
	DISCARD_SUCCESS        = 1
	DISCARD_SMALLER        = 2
	DISCARD_WRONGTURN      = 3
	DISCARD_NOTEXIST       = 4

	BLACKA_GROUP    = 0
	NONBLACKA_GROUP = 1
)

type gameInitData struct {
	startIndex int
}

type dropCards struct {
	PlayerId int
	Cards    pattern.CardPattern
}

type CardGame struct {
	id          int
	rule        rule.ICardRule
	DropList    []dropCards
	Players     []PlayerInfo
	MsgReceiver chan room.GameCommand
	MsgSender   chan room.GameCommand
	End         chan bool
	winGroup    int
	logArea     string
	gameInfo    CardGameInfo
}

func (this CardGame) GetGameId() int {
	return this.id
}
func (this CardGame) SetMsgReceiver(c chan room.GameCommand) {
	this.MsgReceiver = c
}
func (this CardGame) SetMsgSender(c chan room.GameCommand) {
	this.MsgSender = c
}

func (this CardGame) Start() {
	this.rule.Init()
	deal := this.rule.DealCards()
	for i, v := range deal {
		this.Players[i].Cards = v
		if this.hasCards(i, []cards.Card{{CardNumber: 1, CardType: cards.CARDTYPE_CLUB}}) ||
			this.hasCards(i, []cards.Card{{CardNumber: 1, CardType: cards.CARDTYPE_SPADE}}) {
			this.Players[i].Group = BLACKA_GROUP
		} else {
			this.Players[i].Group = NONBLACKA_GROUP
		}
	}
	go this.handleGame()
}

func (this *CardGame) handleGameCommand(pIdx int, cmdString string) {
	var cmd CardCommand
	err := json.Unmarshal([]byte(cmdString), &cmd)
	if err != nil {
		logging.LogError(this.logArea, fmt.Sprintf("invalid card command. %v", cmdString))
	}
	switch cmd.CmdType {
	case CMDTYPE_DISCARD:
		this.Discard(pIdx, cmd.CardList)
	case CMDTYPE_PASS:
		this.Pass(pIdx)
	}
}

func (this *CardGame) handleCommand(gCmd *room.GameCommand) {
	switch gCmd.CmdType {
	case room.GAMECMD_DISCONNECT:
	case room.GAMECMD_RECONNECT:
	case room.GAMECMD_EXIT:
		this.endGame(true)
	case room.GAMECMD_STATUS:
		// notify all
		this.notifyAll()
	case room.GAMECMD_GAME:
		this.handleGameCommand(gCmd.PlayerIndex, gCmd.PlayerCommand)
	}
}

func (this *CardGame) handleGame() {
	logging.LogInfo(this.logArea, fmt.Sprintf("Start to handle game %v", this.id))
GameLoop:
	for {
		select {
		case gCmd := <-this.MsgReceiver:
			this.handleCommand(&gCmd)
		case <-this.End:
			this.endGame(false)
			break GameLoop
		}
	}
	logging.LogInfo(this.logArea, fmt.Sprintf("end to handle game %v", this.id))
}

func (this *CardGame) endGame(isExit bool) {
	this.gameInfo.IsEnd = true
	this.notifyAll()
	var nextData interface{}
	if isExit {
		nextData = nil
	} else {
		nextData = &gameInitData{startIndex: this.gameInfo.Turn}
	}
	this.MsgSender <- room.MakeGameCommandResponse_End(this.gameInfo.ToMessage(), nextData)
}

func (this *CardGame) Clear() {
	for i := range this.Players {
		this.Players[i].Cards = make([]cards.Card, 0)
		this.Players[i].Group = 0
	}
	this.DropList = make([]dropCards, 0)
	this.gameInfo.Wind = false
}

func (this *CardGame) hasCards(playerNumer int, cardList []cards.Card) bool {
	sort.Sort(cards.CardList(this.Players[playerNumer].Cards))
	sort.Sort(cards.CardList(cardList))
	l1, l2 := len(this.Players[playerNumer].Cards), len(cardList)
	for i, j := 0, 0; i < l1 && j < l2; {
		for i < l1 && !this.Players[playerNumer].Cards[i].Equal(&cardList[j]) {
			i++
		}
		if i == l1 {
			return false
		}
		i++
		j++
	}
	return true
}

func (this *CardGame) nextTurn() {
	total := this.rule.PlayerNumber()
	this.gameInfo.Turn = (this.gameInfo.Turn + 1) % total
	for len(this.Players[this.gameInfo.Turn].Cards) == 0 {
		// if there is no one discarding and last turn player has no card, then next one take the turn.
		if this.DropList[len(this.DropList)-1].PlayerId == this.gameInfo.Turn {
			this.gameInfo.Wind = true
		}
		this.gameInfo.Turn = (this.gameInfo.Turn + 1) % total
	}
}

func (this *CardGame) Pass(playerNumber int) bool {
	if this.gameInfo.Turn != playerNumber {
		return false
	}
	ll := len(this.DropList)
	if ll > 0 && (this.DropList[ll-1].PlayerId == playerNumber || this.gameInfo.Wind) {
		return false
	}
	this.nextTurn()
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
	this.DropList = append(this.DropList, dropCards{PlayerId: playerNumber, Cards: pat})

	blackAWin := true
	nonAWin := true
	for _, v := range this.Players {
		if len(v.Cards) != 0 {
			if v.Group == BLACKA_GROUP {
				blackAWin = false
			} else {
				nonAWin = false
			}
		}
	}
	if blackAWin || nonAWin {
		if blackAWin {
			logging.LogInfo_Normal(this.logArea, fmt.Sprintf("black win.\n"))
			this.winGroup = BLACKA_GROUP
		}
		if nonAWin {
			logging.LogInfo_Normal(this.logArea, fmt.Sprintf("non-black win.\n"))
			this.winGroup = NONBLACKA_GROUP
		}
		for _, v := range this.Players {
			if v.Group == this.winGroup {
				v.IsWinner = true
			}
		}
		this.End <- true
	} else {
		this.nextTurn()
	}
}

func (this *CardGame) Discard(playerNumber int, cardList []cards.Card) (int, *pattern.CardPattern) {
	if this.gameInfo.Turn != playerNumber {
		return DISCARD_WRONGTURN, nil
	}
	if !this.hasCards(playerNumber, cardList) {
		return DISCARD_NOTEXIST, nil
	}
	pat := this.rule.GetPatternMatcher()(cardList)
	if pat.PatternType == pattern.PATTERN_INVALID {
		return DISCARD_INVALIDPATTERN, nil
	}
	ll := len(this.DropList)
	if ll == 0 || this.gameInfo.Wind || this.DropList[ll-1].PlayerId == playerNumber {
		this.gameInfo.Wind = false
		this.discard(playerNumber, pat)
		return DISCARD_SUCCESS, &pat
	}
	val, ok := this.rule.Compare(&pat, &this.DropList[ll-1].Cards)
	if ok && val > 0 {
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

func (this CardGame) GetStatus(pIdx int) string {
	//ans := make([]PlayerInfo, len(this.Players))
	//for i, v := range this.Players {
	ans := PlayerInfo{}
	v := this.Players[pIdx]
	ans.Cards = v.Cards
	ans.OnTurn = this.gameInfo.Turn == pIdx

	return ans.ToMessage()
}

func (this *CardGame) GetAllStatusFor(pIdx int) []string {
	ans := make([]string, len(this.Players))
	for i, v := range this.Players {
		p := PlayerInfo{}
		if i == pIdx {
			p.Cards = v.Cards
		} else {
			p.Cards = make([]cards.Card, len(v.Cards))
		}
		p.OnTurn = this.gameInfo.Turn == i
		p.IsWinner = v.IsWinner
		p.Score = v.Score
		ans[i] = p.ToMessage()
	}
	return ans
}

func (this *CardGame) notifyAll() {
	for i, _ := range this.Players {
		this.MsgSender <- room.MakeGameCommandResponse_Notify(i, this.GetAllStatusFor(i), this.gameInfo.ToMessage())
	}
}
