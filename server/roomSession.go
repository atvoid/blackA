package server

import (
	"blackA/room"
)

type RoomSession map[int]*room.Room

var GlobalRoomSession RoomSession = RoomSession{}
