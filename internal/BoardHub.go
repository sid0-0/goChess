package internal

import (
	"gochess/chessBoard"
	"gochess/ws"
)

type ClientType string

const (
	PLAYER_W  ClientType = "player_white"
	PLAYER_B  ClientType = "player_black"
	SPECTATOR ClientType = "spectator"
	LURKER    ClientType = "lurker"
)

type ClientInfoType struct {
	Type ClientType
}

type ClientContextData struct {
	WebsocketClient *ws.Client[ClientInfoType]
	Board           *chessBoard.Board
	Pool            *ws.Pool[ClientInfoType]
}
