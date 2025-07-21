package internal

import (
	"gochess/chessBoard"
	"gochess/ws"
)

type ClientContextData struct {
	Data  *ws.Client
	Board *chessBoard.Board
	Pool  *ws.Pool
}
