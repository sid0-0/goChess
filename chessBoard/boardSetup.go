package chessBoard

import (
	"errors"
	"strings"
	"unicode"
)

func (b *Board) EvaluatePieceMoves() error {
	for ri, squareRow := range b.Squares {
		for ci := range squareRow {
			square := &b.Squares[ri][ci]
			pieceType := square.Piece.PieceType
			if pieceType == PAWN {
				b.loadPawnPieceMoves(square)
			} else if pieceType == KING || pieceType == KNIGHT {
				b.specificSearch(square)
			} else {
				b.fullBoardSearch(square)
			}
		}
	}
	return nil
}

func (b *Board) EvaluateLegalMoves() error {
	b.EvaluatePieceMoves()
	return nil
}

func (b *Board) LoadBoard(fenString string) error {
	verboseFen := ""
	for _, c := range fenString {
		if c >= '1' && c <= '8' {
			verboseFen += strings.Repeat(".", int(c-'0'))
		} else {
			verboseFen += string(c)
		}
	}
	rows := strings.Split(verboseFen, "/")
	if len(rows) != 8 {
		return errors.New("Incorrect row count in FEN")
	}
	for _, row := range rows {
		if len(row) != 8 {
			return errors.New("Incorrect column count in FEN")
		}
	}

	for rindex, row := range rows {
		for cindex, cell := range row {
			// TODO: Store in a temp place then copy only if no errors
			currentSquare := &b.Squares[rindex][cindex]
			if cell == '.' {

			} else if 'a' <= cell && cell <= 'z' {
				currentSquare.Piece.Color = WHITE
			} else if 'A' <= cell && cell <= 'Z' {
				currentSquare.Piece.Color = BLACK
			} else {
				return errors.New("Incorrect characters in FEN[casing]")
			}

			switch unicode.ToLower(cell) {
			case '.':
				currentSquare.Piece.PieceType = EMPTY
			case 'r':
				currentSquare.Piece.PieceType = ROOK
			case 'n':
				currentSquare.Piece.PieceType = KNIGHT
			case 'b':
				currentSquare.Piece.PieceType = BISHOP
			case 'k':
				currentSquare.Piece.PieceType = KING
			case 'q':
				currentSquare.Piece.PieceType = QUEEN
			case 'p':
				currentSquare.Piece.PieceType = PAWN
			default:
				return errors.New("Incorrect characters in FEN[spec]")
			}
		}
	}
	b.EvaluateLegalMoves()
	return nil
}
