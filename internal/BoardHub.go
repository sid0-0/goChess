package internal

import (
	"gochess/chessBoard"
	"gochess/ws"
)

type ClientContextData struct {
	WebsocketClient *ws.Client
	Board           *chessBoard.Board
	Pool            *ws.Pool
}
