package chessBoard

import (
	"fmt"
)

func New() *Board {
	newBoard := Board{}
	newBoard.CastleRights = make(map[COLOR]CastleRight)
	newBoard.CastleRights[WHITE] = CastleRight{}
	newBoard.CastleRights[BLACK] = CastleRight{}
	for i := 0; i < 64; i++ {
		r, c := i/8, i%8
		currentSquare := &newBoard.Squares[r][c]
		currentSquare.Ri = r
		currentSquare.Ci = c
		currentSquare.Rank = fmt.Sprintf("%c", '1'+r)
		currentSquare.File = fmt.Sprintf("%c", 'a'+c)
		if r%2 == c%2 {
			currentSquare.Color = WHITE
		} else {
			currentSquare.Color = BLACK
		}
	}
	// err := newBoard.LoadBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	err := newBoard.LoadBoard("rnb1kbnr/pppppppp/8/3q4/8/4K3/PPPPPPPP/RNBQ1BNR w KQkq - 0 1")
	if err != nil {
		fmt.Println("Error loading board.", err)
	}
	return &newBoard
}

func (b Board) GetRepresentationalSquares() [8][8]Square {
	ans := b.Squares
	for i := 0; i < 4; i++ {
		ans[i], ans[7-i] = ans[7-i], ans[i]
	}
	return ans
}

// func (b *Board) makeMove(oldSquare *Square, newSquare *Square) {
// 	var newData []*Square
// 	for _, diff := range possibleDiff[newSquare.Piece.PieceType] {
// 		nr := newSquare.Ri + diff.r
// 		nc := newSquare.Ci + diff.c
// 		if b.isIndexInRange(nr, nc) && (b.Squares[nr][nc].Piece == nil || newSquare.Piece.Color != b.Squares[nr][nc].Piece.Color) {
// 			newData = append(newData, &b.Squares[nr][nc])
// 		}
// 	}
// 	newSquare.PieceMoves = newData
// }
