package chessBoard

import (
	"slices"
)

func isLegalMove(b *Board, originalFrom *Square, originalTo *Square) bool {
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

	boardCopy.EvaluatePieceMoves()

	kingSquares := boardCopy.findPiece(KING)

	kingSquare := kingSquares[0]
	if kingSquare.Piece.Color != to.Piece.Color {
		kingSquare = kingSquares[1]
	}

	for i := range boardCopy.Squares {
		for j := range boardCopy.Squares[i] {
			square := &boardCopy.Squares[i][j]
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
		// if currentSquare.File == "e" && currentSquare.Rank == "3" && move.File == "d" && move.Rank == "3" {
		// 	fmt.Println("King", currentSquare.File, currentSquare.Rank, move.File, move.Rank, isLegalMove(b, currentSquare, move))
		// }
		if isLegalMove(b, currentSquare, move) {
			currentSquare.LegalMoves = append(currentSquare.LegalMoves, move)
		}
	}
}
