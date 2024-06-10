package chessBoard

import (
	"errors"
	"slices"
)

func (b *Board) MakeMove(oldSquare *Square, newSquare *Square) error {
	if !slices.Contains(oldSquare.LegalMoves, newSquare) {
		return errors.New("invalid move")
	}
	newSquare.Piece = oldSquare.Piece
	oldSquare.Piece = nil

	err := b.EvaluateLegalMoves()
	return err
}
