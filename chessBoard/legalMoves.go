package chessBoard

import (
	"math"
	"slices"
)

// checks if the move is legal by simulating the move and checking if the king is in check
func (b *Board) isLegalMove(originalFrom *Square, originalTo *Square) bool {
	boardCopy := *b

	from := boardCopy.GetSquare(originalFrom.Ri, originalFrom.Ci)
	to := boardCopy.GetSquare(originalTo.Ri, originalTo.Ci)

	if from == nil || to == nil {
		return false
	}
	if from.Piece == nil {
		return true
	}

	to.Piece = from.Piece
	from.Piece = nil

	// validate castling move
	if to.Piece.PieceType == KING && math.Abs(float64(to.Ci-from.Ci)) == 2 {
		pieceColor := to.Piece.Color
		if to.Ci > from.Ci {
			if !boardCopy.CastleRights[pieceColor].Short || !b.isLegalMove(from, b.GetSquare(to.Ri, 5)) {
				return false
			}
		} else {
			if !boardCopy.CastleRights[pieceColor].Long || !b.isLegalMove(from, b.GetSquare(to.Ri, 3)) {
				return false
			}
		}
	}

	boardCopy.EvaluatePieceMoves()

	kingSquares := boardCopy.findPiece(KING)

	kingSquare := kingSquares[0]
	if kingSquare.Piece.Color != to.Piece.Color {
		kingSquare = kingSquares[1]
	}

	for _, row := range boardCopy.Squares {
		for _, square := range row {
			if slices.Contains(square.PieceMoves, kingSquare) {
				return false
			}
		}
	}
	return true
}

func (b *Board) LoadLegalMoves(currentSquare *Square) {
	currentSquare.LegalMoves = []*Square{}
	for _, move := range currentSquare.PieceMoves {
		if b.isLegalMove(currentSquare, move) {
			currentSquare.LegalMoves = append(currentSquare.LegalMoves, move)
		}
	}
}
