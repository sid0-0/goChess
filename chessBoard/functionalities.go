package chessBoard

import (
	"errors"
	"slices"
)

func (b *Board) MakeMove(oldSquare *Square, newSquare *Square) error {
	if oldSquare.Piece.Color != b.Turn {
		return errors.New("not your turn")
	}
	if !slices.Contains(oldSquare.LegalMoves, newSquare) {
		return errors.New("invalid move")
	}
	newSquare.Piece = oldSquare.Piece
	oldSquare.Piece = nil

	err := b.EvaluateLegalMoves()
	if err != nil {
		return err
	}

	if b.Turn == BLACK {
		b.Turn = WHITE
	} else {
		b.Turn = BLACK
	}
	return nil
}
