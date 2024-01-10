package chessBoard

func isIndexInRange(b *Board, r int, c int) bool {
	if r <= 0 || c <= 0 {
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
