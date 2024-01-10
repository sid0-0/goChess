package chessBoard

import (
	"fmt"
)

func New() *Board {
	newBoard := Board{}
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
		err := newBoard.LoadBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
		if err != nil {
			fmt.Println("Error loading board.", err)
		}
	}
	return &newBoard
}

func (b Board) GetRepresentationalSquares() [8][8]Square {
	ans := b.Squares
	var temp [8]Square
	for i := 0; i < 4; i++ {
		temp = ans[i]
		ans[i] = ans[7-i]
		ans[7-i] = temp
	}
	return ans
}
