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
	NONE_GROUP      = 2
)

type gameInitData struct {
	startIndex int
}

type dropCards struct {
	PlayerId int
	Cards    pattern.CardPattern
}

type CardGame struct {
	id               int
	rule             rule.ICardRule
	DropList         []dropCards
	Players          []PlayerInfo
	MsgReceiver      chan room.GameCommand
	MsgSender        chan room.GameCommand
	End              chan bool
	logArea          string
	gameInfo         CardGameInfo
	firstFinishGroup int
}

func (this *CardGame) GetGameId() int {
	return this.id
}
func (this *CardGame) SetMsgReceiver(c chan room.GameCommand) {
	this.MsgReceiver = c
}
func (this *CardGame) SetMsgSender(c chan room.GameCommand) {
	this.MsgSender = c
}

func (this *CardGame) Start() {
	this.rule.Init()
	deal := this.rule.DealCards()
	this.firstFinishGroup = -1
	for i, v := range deal {
		this.Players[i].Cards = v
		this.Players[i].IsClear = true
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
		a, _ := this.Discard(pIdx, cmd.CardList)
		if a == DISCARD_SUCCESS {
			logging.LogInfo_Detail(this.logArea, fmt.Sprintf("index %v discarded successfully.", pIdx))
			this.notifyCommand(cmd)
		} else {
			logging.LogInfo_Detail(this.logArea, fmt.Sprintf("index %v failed to discard, ErrorCode: %v.", pIdx, a))
		}
		this.notifyAll()
	case CMDTYPE_PASS:
		ok := this.Pass(pIdx)
		if ok {
			logging.LogInfo_Detail(this.logArea, fmt.Sprintf("index %v passed successfully.", pIdx))
			this.notifyCommand(cmd)
		} else {
			logging.LogInfo_Detail(this.logArea, fmt.Sprintf("index %v failed to pass.", pIdx))
		}
		this.notifyAll()
	}
}

func (this *CardGame) handleCommand(gCmd *room.GameCommand) {
	switch gCmd.CmdType {
	case room.GAMECMD_DISCONNECT:
		this.notifyAll()
	case room.GAMECMD_RECONNECT:
		this.notifyAll()
	case room.GAMECMD_EXIT:
		this.End <- true
	case room.GAMECMD_STATUS:
		// notify all
		this.notifyAll()
	case room.GAMECMD_GAME:
		this.handleGameCommand(gCmd.PlayerIndex, gCmd.PlayerCommand)
	}
}

func (this *CardGame) handleGame() {
	logging.LogInfo(this.logArea, fmt.Sprintf("Start to handle game %v", this.id))
	this.notifyAll()
GameLoop:
	for {
		select {
		case gCmd := <-this.MsgReceiver:
			this.handleCommand(&gCmd)
		case isExit := <-this.End:
			this.endGame(isExit)
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
	logging.LogInfo(this.logArea, fmt.Sprintf("game finished. %v", isExit))
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
		this.Players[this.gameInfo.Turn].DropCards = nil
		this.Players[this.gameInfo.Turn].IsClear = true
		this.gameInfo.Turn = (this.gameInfo.Turn + 1) % total
	}
	this.Players[this.gameInfo.Turn].DropCards = nil
	this.Players[this.gameInfo.Turn].IsClear = true

	// if wind or no one bid, turn to next round discarding
	if this.gameInfo.Wind || this.gameInfo.Turn == this.DropList[len(this.DropList)-1].PlayerId {
		logging.LogInfo_Detail(this.logArea, fmt.Sprintf("start next round. Turn:%v", this.gameInfo.Turn))
		for i, _ := range this.Players {
			this.Players[i].DropCards = nil
			this.Players[i].IsClear = true
		}
	}
}

func (this *CardGame) Pass(playerNumber int) bool {
	if this.gameInfo.Turn != playerNumber {
		return false
	}
	ll := len(this.DropList)
	if ll == 0 {
		return false
	}
	if ll > 0 && (this.DropList[ll-1].PlayerId == playerNumber || this.gameInfo.Wind) {
		return false
	}

	this.Players[playerNumber].DropCards = make([]cards.Card, 0)
	this.Players[playerNumber].IsClear = false

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

		if this.firstFinishGroup == -1 && len(this.Players[playerNumber].Cards) == 0 {
			this.firstFinishGroup = this.Players[playerNumber].Group
		}
	}
	this.DropList = append(this.DropList, dropCards{PlayerId: playerNumber, Cards: pat})

	this.Players[playerNumber].DropCards = pat.CardList
	this.Players[playerNumber].IsClear = false

	// judge whether anyone wins
	blackAFinish := true
	nonAFinish := true
	for _, v := range this.Players {
		if len(v.Cards) != 0 {
			if v.Group == BLACKA_GROUP {
				blackAFinish = false
			} else {
				nonAFinish = false
			}
		}
	}
	if blackAFinish || nonAFinish {
		if blackAFinish && this.firstFinishGroup == BLACKA_GROUP {
			logging.LogInfo_Normal(this.logArea, fmt.Sprintf("black win.\n"))
			this.gameInfo.WinGroup = BLACKA_GROUP
		} else if nonAFinish && this.firstFinishGroup == NONBLACKA_GROUP {
			logging.LogInfo_Normal(this.logArea, fmt.Sprintf("non-black win.\n"))
			this.gameInfo.WinGroup = NONBLACKA_GROUP
		} else {
			this.gameInfo.WinGroup = NONE_GROUP
			logging.LogInfo_Normal(this.logArea, fmt.Sprintf("draw.\n"))
		}
		for i, v := range this.Players {
			if v.Group == this.gameInfo.WinGroup {
				this.Players[i].IsWinner = true
			}
		}
		this.End <- false
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
		if v.Group == this.gameInfo.WinGroup {
			winner = append(winner, i)
		}
	}
	return winner
}

func (this *CardGame) GetStatus(pIdx int) string {
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
		p.DropCards = v.DropCards
		p.IsClear = v.IsClear
		p.IsWinner = v.IsWinner
		ans[i] = p.ToMessage()
	}
	return ans
}

func (this *CardGame) notifyAll() {
	logging.LogInfo_Detail(this.logArea, "notify all status.")
	for i, _ := range this.Players {
		this.MsgSender <- room.MakeGameCommandResponse_Notify(i, this.GetAllStatusFor(i), this.gameInfo.ToMessage())
	}
}

func (this *CardGame) notifyCommand(cCmd CardCommand) {
	for i, _ := range this.Players {
		this.MsgSender <- room.MakeGameCommandResponse_Command(i, cCmd.ToMessage())
	}
}
