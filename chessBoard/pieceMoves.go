package chessBoard

type PositionalDiff struct {
	r int
	c int
}

var possibleDiff = map[PIECE_TYPE][]PositionalDiff{
	// evaluate
	KNIGHT: {
		PositionalDiff{1, 2},
		PositionalDiff{2, 1},
		PositionalDiff{2, -1},
		PositionalDiff{1, -2},
		PositionalDiff{-1, -2},
		PositionalDiff{-2, -1},
		PositionalDiff{-2, 1},
		PositionalDiff{-2, 2},
	},
	KING: {
		PositionalDiff{1, 0},
		PositionalDiff{-1, 0},
		PositionalDiff{0, 1},
		PositionalDiff{0, -1},
		PositionalDiff{1, 1},
		PositionalDiff{-1, -1},
		PositionalDiff{1, -1},
		PositionalDiff{-1, 1},
	},
	// Pawn assumes a white pawn's diff
	PAWN: {
		PositionalDiff{1, 0},
		PositionalDiff{1, 1},
		PositionalDiff{1, -1},
	},
	// search and evaluate
	ROOK: {
		PositionalDiff{1, 0},
		PositionalDiff{-1, 0},
		PositionalDiff{0, 1},
		PositionalDiff{0, -1},
	},
	BISHOP: {
		PositionalDiff{1, 1},
		PositionalDiff{-1, -1},
		PositionalDiff{1, -1},
		PositionalDiff{-1, 1},
	},
	QUEEN: {
		PositionalDiff{1, 0},
		PositionalDiff{-1, 0},
		PositionalDiff{0, 1},
		PositionalDiff{0, -1},
		PositionalDiff{1, 1},
		PositionalDiff{-1, -1},
		PositionalDiff{1, -1},
		PositionalDiff{-1, 1},
	},
}

func (b Board) fullBoardSearch(currentSquare *Square) {
	var newData []*Square
	flag := 0
	for i := 0; i < 8; i++ {
		flag = 0
		for _, diff := range possibleDiff[currentSquare.Piece.PieceType] {
			nr, nc := currentSquare.Ri+(i*diff.r), currentSquare.Ci+(i*diff.c)
			if 0 <= nr && nr < 8 && 0 <= nc && nc < 8 {
				newData = append(newData, &b.Squares[nr][nc])
				flag = 1
			}
		}
		if flag == 0 {
			break
		}
	}
	currentSquare.PieceMoves = newData
}

func (b Board) specificSearch(currentSquare *Square) {
	var newData []*Square
	for _, diff := range possibleDiff[currentSquare.Piece.PieceType] {
		nr := currentSquare.Ri + diff.r
		nc := currentSquare.Ci + diff.c
		if 0 <= nr && nr < 8 && 0 <= nc && nc < 8 {
			newData = append(newData, &b.Squares[nr][nc])
		}
	}
	currentSquare.PieceMoves = newData
}

func (b Board) loadPawnPieceMoves(currentSquare *Square) {
	var newData []*Square
	var updatedRow int
	if currentSquare.Color == WHITE {
		updatedRow = currentSquare.Ri + 1
	} else {
		updatedRow = currentSquare.Ri - 1
	}

	if (currentSquare.Color == WHITE && updatedRow < 7) || (currentSquare.Color == BLACK && updatedRow > 0) {
		newData = append(newData, &b.Squares[updatedRow][currentSquare.Ci])
		if currentSquare.Ci > 0 {
			newData = append(newData, &b.Squares[updatedRow][currentSquare.Ci-1])
		}
		if currentSquare.Ci < 7 {
			newData = append(newData, &b.Squares[updatedRow][currentSquare.Ci+1])
		}
	}
	currentSquare.PieceMoves = newData
}

func (b Board) loadRookPieceMoves(row int, column int) {
	currentSquare := &b.Squares[row][column]
	var newData []*Square
	for i := 0; i < 8; i++ {
		if i != row {
			newData = append(newData, &b.Squares[i][column])
		}
		if i != column {
			newData = append(newData, &b.Squares[row][i])
		}
	}
	currentSquare.PieceMoves = newData
}
