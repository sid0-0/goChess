package chessBoard

func (b *Board) GetSquare(r int, c int) *Square {
	if !b.isIndexInRange(r, c) {
		return nil
	}
	return &b.Squares[r][c]
}

func (b *Board) findPiece(pt PIECE_TYPE) []*Square {
	var ans []*Square
	for i := range b.Squares {
		for j := range b.Squares[i] {
			square := &b.Squares[i][j]
			if square.Piece != nil && square.Piece.PieceType == pt {
				ans = append(ans, square)
			}
		}
	}
	return ans
}

func (b *Board) isIndexInRange(r int, c int) bool {
	if r < 0 || c < 0 || b == nil {
		return false
	}
	rc := len(b.Squares)
	if rc <= r {
		return false
	}
	cc := len(b.Squares[0])
	if cc <= c {
		return false
	}
	return true
}

func (b *Board) hasNoMoves() bool {
	turnColor := b.Turn
	for _, row := range b.Squares {
		for _, square := range row {
			if square.Piece != nil && square.Piece.Color == turnColor && len(square.LegalMoves) > 0 {
				return false
			}
		}
	}
	return true
}
