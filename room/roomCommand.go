package room

import (
	"blackA/user"
	"encoding/json"
)

const (
	// room command for user to room
	ROOMCMD_JOIN       = 0
	ROOMCMD_LEAVE      = 1
	ROOMCMD_DISCONNECT = 2
	ROOMCMD_RECONNECT  = 3
	ROOMCMD_READY      = 4
	ROOMCMD_NOTREADY   = 5
	ROOMCMD_ROOMSTATUS = 6 // room status contains room member information
	ROOMCMD_GAMESTATUS = 7 // game statuc contains room status information, game info and player info

	ROOMCMD_RESPONSE_JOIN_SUCCESS = 100
	ROOMCMD_RESPONSE_JOIN_FULL    = 101
	ROOMCMD_RESPONSE_JOIN_STARTED = 102

	ROOMCMD_RESPONSE_LEAVE_SUCCESS = 200

	ROOMCMD_RESPONSE_RECONNECT_SUCCESS = 400

	ROOMCMD_RESPONSE_ROOMSTATUS = 500
	ROOMCMD_RESPONSE_GAMESTATUS = 600
)

type RoomCommand struct {
	RoomId      int
	CmdType     int
	GameStarted bool
	UserInfo    []*UserInfo
	GameInfo    string
}

func (this *RoomCommand) ToMessage() string {
	b, _ := json.Marshal(*this)
	return string(b)
}

func MakeRoomRequest_Reconnect(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RECONNECT}
}

func MakeRoomRequest_Disconnect(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_DISCONNECT}
}

func MakeRoomResponse_Join_Success(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_JOIN_SUCCESS}
}

func MakeRoomResponse_Join_Full(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_JOIN_FULL}
}

func MakeRoomResponse_Join_Started(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_JOIN_STARTED}
}

func MakeRoomResponse_Leave_Success(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_LEAVE_SUCCESS}
}

func MakeRoomResponse_Reconnect_SUCCESS(rId int) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_RECONNECT_SUCCESS}
}

func MakeRoomResponse_RoomStatus(rId int, uInfo []*UserInfo, started bool) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_ROOMSTATUS, GameStarted: started, UserInfo: uInfo}
}

func MakeRoomResponse_GameStatus(rId int, uInfo []*UserInfo, gameInfo string, started bool) RoomCommand {
	return RoomCommand{RoomId: rId, CmdType: ROOMCMD_RESPONSE_GAMESTATUS, GameInfo: gameInfo, GameStarted: started, UserInfo: uInfo}
}

func MakeUserCommandForRoom(uid int, rCmd RoomCommand) user.Command {
	return user.Command{CmdType: user.CMDTYPE_ROOM, UserId: uid, Command: rCmd.ToMessage()}
}
