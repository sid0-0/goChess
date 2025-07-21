package internal

import (
	"gochess/chessBoard"
	"gochess/ws"
)

type ClientContextData struct {
	WebSocketData *ws.Client
	Board         *chessBoard.Board
	Pool          *ws.Pool
}
