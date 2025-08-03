package chessBoard

import (
	"errors"
	"math"
	"slices"
)

type MoveArgs struct {
	FromSquare         *Square
	ToSquare           *Square
	PromotionPieceType PIECE_TYPE
}

func (b *Board) IsPromotionMove(args MoveArgs) bool {
	from := args.FromSquare
	to := args.ToSquare
	if from == nil || to == nil {
		return false
	}
	if from.Piece.PieceType == PAWN && (to.Ri == 0 || to.Ri == 7) {
		return true
	}
	return false
}

func (b *Board) MakeMove(args MoveArgs) error {
	oldSquare := args.FromSquare
	newSquare := args.ToSquare
	promotionPieceType := args.PromotionPieceType
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
		// Pawn moved 2 squares forward
		if math.Abs(float64(oldSquare.Ri-newSquare.Ri)) == 2 {
			b.EnPassantSquare = b.GetSquare((oldSquare.Ri+newSquare.Ri)/2, oldSquare.Ci)
		}
		if b.EnPassantSquare == newSquare {
			b.EnPassantSquare = nil
			squareToClear := b.GetSquare(oldSquare.Ri, newSquare.Ci)
			squareToClear.Piece = nil
		}
		// Handle pawn promotion
		if b.IsPromotionMove(args) {
			newSquare.Piece.PieceType = promotionPieceType
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

func (b *Board) IsCheck() bool {
	turnColor := b.Turn
	kingSquareCandidates := b.findPiece(KING)
	kingSquare := kingSquareCandidates[0]
	if kingSquare.Piece.Color != turnColor {
		kingSquare = kingSquareCandidates[1]
	}
	for _, row := range b.Squares {
		for _, square := range row {
			if square.Piece == nil || square.Piece.Color == turnColor {
				continue
			}
			if slices.Contains(square.LegalMoves, kingSquare) {
				return true
			}
		}
	}
	return false
}

func (b *Board) IsStalemate() bool {
	return b.hasNoMoves() && !b.IsCheck()
}

func (b *Board) IsCheckmate() bool {
	return b.hasNoMoves() && b.IsCheck()
}
