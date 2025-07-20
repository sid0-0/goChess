package internal

import (
	"gochess/chessBoard"
	"gochess/ws"
)

type BoardHub struct {
	Board *chessBoard.Board
	Hub   *ws.Hub
}

func NewBoardHub() *BoardHub {
	return &BoardHub{
		Board: chessBoard.New(),
		Hub:   ws.NewHub(),
	}
}
