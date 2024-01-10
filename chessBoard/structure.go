package chessBoard

type COLOR int

const (
	EMPTY_COLOR COLOR = iota
	WHITE
	BLACK
)

type PIECE_TYPE string

const (
	EMPTY  PIECE_TYPE = " "
	PAWN   PIECE_TYPE = "P"
	KNIGHT PIECE_TYPE = "N"
	BISHOP PIECE_TYPE = "B"
	ROOK   PIECE_TYPE = "R"
	QUEEN  PIECE_TYPE = "Q"
	KING   PIECE_TYPE = "K"
)

type Piece struct {
	Color     COLOR
	PieceType PIECE_TYPE
}

type Square struct {
	Ci         int
	Ri         int
	File       string
	Rank       string
	Color      COLOR
	LegalMoves []*Square
	PieceMoves []*Square
	Piece      Piece
}

type CastleRight struct {
	Long  bool
	Short bool
}

type Board struct {
	Squares         [8][8]Square
	Turn            COLOR
	CastleRights    map[COLOR]CastleRight
	EnPassantSquare *Square
	HalfMoveCounter int
	FullMoveCounter int
}
