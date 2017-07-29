package room

import (
	"encoding/json"
)

const (
	GAMECMD_RESPONSE_NOTIFY  = 0
	GAMECMD_RESPONSE_COMMAND = 1
	GAMECMD_RESPONSE_GAMEEND = 2

	GAMECMD_GAME       = 100
	GAMECMD_RECONNECT  = 101
	GAMECMD_DISCONNECT = 102
	GAMECMD_EXIT       = 103
	GAMECMD_STATUS     = 104
)

type IGame interface {
	GetGameId() int
	GetStatus(pIdx int) string
	SetMsgReceiver(chan GameCommand)
	SetMsgSender(chan GameCommand)
	Start()
}

type IGameCreator interface {
	CreateGame(initData interface{}) IGame
	GetPlayerCapacity() int
}

type GameCommand struct {
	CmdType       int
	PlayerInfo    []string
	PlayerIndex   int
	PlayerCommand string
	GameInfo      string
	NextGameData  interface{}
}

func (this *GameCommand) ToMessage() string {
	b, _ := json.Marshal(*this)
	return string(b)
}

func MakeGameCommandRequest_Game(pIdx int, pCommand string) GameCommand {
	return GameCommand{CmdType: GAMECMD_GAME, PlayerCommand: pCommand, PlayerIndex: pIdx}
}

func MakeGameCommandRequest_Reconnect(pIdx int) GameCommand {
	return GameCommand{CmdType: GAMECMD_RECONNECT, PlayerIndex: pIdx}
}

func MakeGameCommandRequest_Disconnect(pIdx int) GameCommand {
	return GameCommand{CmdType: GAMECMD_DISCONNECT, PlayerIndex: pIdx}
}

func MakeGameCommandRequest_Exit(pIdx int) GameCommand {
	return GameCommand{CmdType: GAMECMD_EXIT, PlayerIndex: pIdx}
}

func MakeGameCommandRequest_Status(pIdx int) GameCommand {
	return GameCommand{CmdType: GAMECMD_STATUS, PlayerIndex: pIdx}
}

func MakeGameCommandResponse_Notify(pIdx int, playInfos []string, gInfo string) GameCommand {
	return GameCommand{CmdType: GAMECMD_RESPONSE_NOTIFY, PlayerInfo: playInfos, GameInfo: gInfo, PlayerIndex: pIdx}
}

func MakeGameCommandResponse_Command(pIdx int, pCmd string) GameCommand {
	return GameCommand{CmdType: GAMECMD_RESPONSE_COMMAND, PlayerCommand: pCmd, PlayerIndex: pIdx}
}

// this command is only for ROOM only, will not expose to players
func MakeGameCommandResponse_End(gInfo string, nextGameData interface{}) GameCommand {
	return GameCommand{CmdType: GAMECMD_RESPONSE_GAMEEND, GameInfo: gInfo, NextGameData: nextGameData}
}
