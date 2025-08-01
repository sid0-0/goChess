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
		PositionalDiff{-1, 2},
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

func (b *Board) fullBoardSearch(currentSquare *Square) {
	var newData []*Square
	for _, diff := range possibleDiff[currentSquare.Piece.PieceType] {
		i := 1
		for {
			nr, nc := currentSquare.Ri+(i*diff.r), currentSquare.Ci+(i*diff.c)
			if b.isIndexInRange(nr, nc) && (b.Squares[nr][nc].Piece == nil || currentSquare.Piece.Color != b.Squares[nr][nc].Piece.Color) {
				newData = append(newData, &b.Squares[nr][nc])
				if b.Squares[nr][nc].Piece != nil && currentSquare.Piece.Color != b.Squares[nr][nc].Piece.Color {
					break
				}
			} else {
				break
			}
			i++
		}
	}
	currentSquare.PieceMoves = newData
}

func (b *Board) specificSearch(currentSquare *Square) {
	var newData []*Square
	for _, diff := range possibleDiff[currentSquare.Piece.PieceType] {
		nr := currentSquare.Ri + diff.r
		nc := currentSquare.Ci + diff.c
		if b.isIndexInRange(nr, nc) && (b.Squares[nr][nc].Piece == nil || currentSquare.Piece.Color != b.Squares[nr][nc].Piece.Color) {
			newData = append(newData, &b.Squares[nr][nc])
		}
	}
	currentSquare.PieceMoves = newData
}

func (b *Board) loadKingPieceMoves(currentSquare *Square) {
	// load direct movement
	b.specificSearch(currentSquare)
	// load castling moves
	if currentSquare.Piece.Color == WHITE {
		if b.CastleRights[WHITE].Long {
			if b.Squares[0][3].Piece == nil && b.Squares[0][2].Piece == nil && b.Squares[0][1].Piece == nil {
				currentSquare.PieceMoves = append(currentSquare.PieceMoves, b.GetSquare(0, 2))
			}
		}
		if b.CastleRights[WHITE].Short {
			if b.Squares[0][5].Piece == nil && b.Squares[0][6].Piece == nil {
				currentSquare.PieceMoves = append(currentSquare.PieceMoves, b.GetSquare(0, 6))
			}
		}
	} else {
		if b.CastleRights[BLACK].Long {
			if b.Squares[7][3].Piece == nil && b.Squares[7][2].Piece == nil && b.Squares[7][1].Piece == nil {
				currentSquare.PieceMoves = append(currentSquare.PieceMoves, b.GetSquare(7, 2))
			}
		}
		if b.CastleRights[BLACK].Short {
			if b.Squares[7][5].Piece == nil && b.Squares[7][6].Piece == nil {
				currentSquare.PieceMoves = append(currentSquare.PieceMoves, b.GetSquare(7, 6))
			}
		}
	}

	pieceMoveSquares := make([]string, len(currentSquare.PieceMoves))
	for i, square := range currentSquare.PieceMoves {
		pieceMoveSquares[i] = square.File + square.Rank
	}
	// spew.Println(currentSquare.Piece.PieceType+" Moves for", currentSquare.Piece.Color, ":", b.GetSquare(0, 2).Ri, b.GetSquare(0, 2).Ci)
}

func (b *Board) loadPawnPieceMoves(currentSquare *Square) {
	var newData []*Square
	multiplier := 1
	if currentSquare.Piece.Color == BLACK {
		multiplier = -1
	}

	// natural movement
	movementDiff := []PositionalDiff{{1, 0}}
	// 2 block movement only on init square
	if (currentSquare.Piece.Color == BLACK && currentSquare.Ri == 6) || (currentSquare.Piece.Color == WHITE && currentSquare.Ri == 1) {
		movementDiff = append(movementDiff, PositionalDiff{2, 0})
	}

	for _, diff := range movementDiff {
		nr, nc := currentSquare.Ri+(multiplier*diff.r), currentSquare.Ci+(multiplier*diff.c)
		if b.isIndexInRange(nr, nc) && b.Squares[nr][nc].Piece == nil {
			newData = append(newData, &b.Squares[nr][nc])
		} else {
			break
		}
	}

	// check attacking movements
	for _, diff := range []PositionalDiff{{1, 1}, {1, -1}} {
		nr := currentSquare.Ri + multiplier*diff.r
		nc := currentSquare.Ci + multiplier*diff.c
		targetSquare := b.GetSquare(nr, nc)
		if targetSquare == nil {
			continue
		}
		isAttackingAPiece := targetSquare.Piece != nil && currentSquare.Piece.Color != targetSquare.Piece.Color
		isEnPassantPossible := (targetSquare == b.EnPassantSquare)
		if isAttackingAPiece || isEnPassantPossible {
			newData = append(newData, targetSquare)
		}
	}

	currentSquare.PieceMoves = newData
}
