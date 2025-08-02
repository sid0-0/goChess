package chessBoard

import (
	"errors"
	"math"
	"slices"
)

func (b *Board) MakeMove(oldSquare *Square, newSquare *Square) error {
	if oldSquare == nil || newSquare == nil || oldSquare == newSquare {
		return errors.New("invalid squares")
	}
	if oldSquare.Piece.Color != b.Turn {
		return errors.New("not your turn")
	}
	if !slices.Contains(oldSquare.LegalMoves, newSquare) {
		return errors.New("invalid move")
	}
	newSquare.Piece = oldSquare.Piece
	oldSquare.Piece = nil

	if newSquare.Piece.PieceType == PAWN {
		if math.Abs(float64(oldSquare.Ri-newSquare.Ri)) == 2 {
			b.EnPassantSquare = b.GetSquare((oldSquare.Ri+newSquare.Ri)/2, oldSquare.Ci)
		}
		if b.EnPassantSquare == newSquare {
			b.EnPassantSquare = nil
			squareToClear := b.GetSquare(oldSquare.Ri, newSquare.Ci)
			squareToClear.Piece = nil
		}
	}

	if newSquare.Piece.PieceType == KING {
		b.CastleRights[newSquare.Piece.Color].Long = false
		b.CastleRights[newSquare.Piece.Color].Short = false
		// If the king moves two squares, it is a castling move
		if math.Abs(float64(oldSquare.Ci-newSquare.Ci)) == 2 {
			// Short castling col indices
			oldRookSquareCol, newRookSquareCol := 0, 3
			// Long castling cold indices
			if newSquare.Ci > oldSquare.Ci {
				oldRookSquareCol, newRookSquareCol = 7, 5
			}
			rookSquare := b.GetSquare(newSquare.Ri, oldRookSquareCol)
			newRookSquare := b.GetSquare(newSquare.Ri, newRookSquareCol)
			newRookSquare.Piece = rookSquare.Piece
			rookSquare.Piece = nil
		}
	}

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
