package chessBoard

import "slices"

func isLegalMove(b *Board, originalFrom *Square, originalTo *Square) bool {
	boardCopy := *b

	from := boardCopy.GetSquare(originalFrom.Ri, originalFrom.Ci)
	to := boardCopy.GetSquare(originalTo.Ri, originalTo.Ci)
	if from == nil || to == nil {
		return false
	}

	to.Piece = from.Piece
	from.Piece = nil

	boardCopy.EvaluatePieceMoves()

	kingSquare := boardCopy.findPiece(KING)[0]
	for i := range b.Squares {
		for j := range b.Squares[i] {
			square := &b.Squares[i][j]
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
		if isLegalMove(b, currentSquare, move) {
			currentSquare.LegalMoves = append(currentSquare.LegalMoves, move)
		}
	}
}
