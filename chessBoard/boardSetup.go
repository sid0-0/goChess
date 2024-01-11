package chessBoard

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func (b *Board) EvaluatePieceMoves() error {
	for i := range b.Squares {
		for j := range b.Squares[i] {
			square := &b.Squares[i][j]
			if square.Piece == nil {
				continue
			}
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
	for i := range b.Squares {
		for j := range b.Squares[i] {
			b.LoadLegalMoves(&b.Squares[i][j])
		}
	}
	return nil
}

func (b *Board) LoadBoard(fenString string) error {
	// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	parts := strings.Split(fenString, " ")

	// Assessing Turn
	if parts[1] == "w" {
		b.Turn = WHITE
	} else {
		b.Turn = BLACK
	}

	// Assessing CastleRights
	newCastleRights := map[COLOR]CastleRight{
		WHITE: {},
		BLACK: {},
	}
	if rights, ok := b.CastleRights[WHITE]; ok {
		if strings.ContainsRune(parts[2], 'K') {
			rights.Long = true
		}
		if strings.ContainsRune(parts[2], 'Q') {
			rights.Short = true
		}
	}

	if rights, ok := b.CastleRights[BLACK]; ok {
		if strings.ContainsRune(parts[2], 'k') {
			rights.Long = true
		}
		if strings.ContainsRune(parts[2], 'q') {
			rights.Short = true
		}
	}

	b.CastleRights = newCastleRights

	// Assessing EnPassantSquare
	if len(parts[3]) > 2 {
		return errors.New("Incorrect EnPassant Square in FEN")
	}
	if parts[3] != "-" {
		r, c := int(parts[3][0]-'a'), int(parts[3][1]-'1')
		b.EnPassantSquare = &b.Squares[r][c]
	}

	// Loading HalfMoveCounter and FullMoveCounter
	val, err := strconv.Atoi(parts[4])
	if err != nil {
		return errors.New("Incorrect Half move count format")
	}
	b.HalfMoveCounter = val

	val, err = strconv.Atoi(parts[4])
	if err != nil {
		return errors.New("Incorrect Full move count format")
	}
	b.FullMoveCounter = val

	// Generating Board configuration
	if len(parts) != 1 && len(parts) != 6 {
		return errors.New("Insufficient data in FEN")
	}

	placement := ""
	for _, c := range parts[0] {
		if c >= '1' && c <= '8' {
			placement += strings.Repeat(".", int(c-'0'))
		} else {
			placement += string(c)
		}
	}
	rows := strings.Split(placement, "/")
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
			currentSquare := &b.Squares[7-rindex][cindex]
			if cell != '.' {
				currentSquare.Piece = &Piece{}
			}
			if cell == '.' {

			} else if 'a' <= cell && cell <= 'z' {
				currentSquare.Piece.Color = BLACK
			} else if 'A' <= cell && cell <= 'Z' {
				currentSquare.Piece.Color = WHITE
			} else {
				return errors.New("Incorrect characters in FEN[casing]")
			}

			switch unicode.ToLower(cell) {
			case '.':

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
	// b.Squares[0][1].PieceMoves = append(b.Squares[0][1].PieceMoves, &b.Squares[2][2])
	return nil
}
