package cardgame

import (
	"blackA/cards/rule"
	"blackA/room"
	"fmt"
)

var gid = 0

type CardGameCreator struct {
}

func (this *CardGameCreator) CreateGame(initData interface{}) room.IGame {
	gid++
	startTurn := 0
	if initData != nil {
		data, ok := initData.(gameInitData)
		if ok {
			startTurn = data.startIndex
		}
	}
	game := CardGame{
		id:       gid,
		rule:     rule.ICardRule(&rule.BlackAGame{}),
		gameInfo: CardGameInfo{Turn: startTurn, Wind: false, IsEnd: false},
	}
	game.Players = make([]PlayerInfo, game.rule.PlayerNumber())
	game.DropList = make([]dropCards, 0, 60)
	game.End = make(chan bool)
	game.gameInfo.Wind = false
	game.logArea = area + fmt.Sprintf("_%v", gid)
	return game
}

func (this *CardGameCreator) GetPlayerCapacity() int {
	return 4
}
