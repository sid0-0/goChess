package internal

import (
	"encoding/json"
	"errors"
	"gochess/chessBoard"
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig/v3"
)

func LoadAllTemplates(location string) (*template.Template, error) {
	var allTemplates *template.Template
	log.Println("Loading templates from:", location)
	allTemplates = template.New("").Funcs(sprig.FuncMap())
	err := filepath.WalkDir(location, func(path string, _ fs.DirEntry, _ error) error {
		if strings.HasSuffix(path, ".html") {
			_, err := allTemplates.ParseFiles(path)
			if err == nil {
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}

func GetLoadLegalMovesJson(board *chessBoard.Board) template.JS {
	if board != nil {
		allLegalMoves := map[string][]string{}
		// collect all legal moves in an object
		for _, row := range board.Squares {
			for _, square := range row {
				// if square.Piece == nil || square.Piece.Color != board.Turn {
				// 	continue // skip squares that are not the current player's turn
				// }
				squareNotation := square.File + square.Rank
				legalMovesForSquare := []string{}
				for _, square := range square.LegalMoves {
					legalMovesForSquare = append(legalMovesForSquare, square.File+square.Rank)
				}
				allLegalMoves[squareNotation] = legalMovesForSquare
			}
		}
		dataMapJson, err := json.Marshal(allLegalMoves)

		if err == nil {
			return template.JS(dataMapJson)
		}
	}
	return template.JS("")
}

func ResolveSquare(board *chessBoard.Board, squareId string) *chessBoard.Square {
	if len(squareId) != 2 {
		return nil
	}
	ri, ci := int(squareId[1]-'1'), int(squareId[0]-'a')
	square := board.GetSquare(ri, ci)
	return square
}

type MakeMoveArgs struct {
	Board              *chessBoard.Board
	PlayerType         ClientType
	FromSquareId       string
	ToSquareId         string
	PromotionPieceType chessBoard.PIECE_TYPE
}

func ResolveSquareAndCheckPromotion(args MakeMoveArgs) (bool, error) {
	board := args.Board
	fromSquareId, toSquareId := args.FromSquareId, args.ToSquareId

	fromSquare := ResolveSquare(board, fromSquareId)
	toSquare := ResolveSquare(board, toSquareId)

	if toSquare == nil || fromSquare == nil {
		return false, errors.New("invalid square")
	}

	return board.IsPromotionMove(chessBoard.MoveArgs{
		FromSquare: fromSquare,
		ToSquare:   toSquare,
	}), nil
}

func ResolveSquareAndMakeMove(args MakeMoveArgs) error {
	board := args.Board
	playerType := args.PlayerType
	fromSquareId, toSquareId := args.FromSquareId, args.ToSquareId
	promotionPieceType := args.PromotionPieceType

	if !((playerType == PLAYER_W && board.Turn == chessBoard.WHITE) || (playerType == PLAYER_B && board.Turn == chessBoard.BLACK)) {
		return errors.New("it's not your turn")
	}
	fromSquare := ResolveSquare(board, fromSquareId)
	toSquare := ResolveSquare(board, toSquareId)

	if toSquare == nil || fromSquare == nil {
		return errors.New("invalid square")
	}
	err := board.MakeMove(chessBoard.MoveArgs{
		FromSquare:         fromSquare,
		ToSquare:           toSquare,
		PromotionPieceType: promotionPieceType,
	})
	return err
}

func GetBoardPlayerColorFromPlayerType(playerType ClientType) chessBoard.COLOR {
	if playerType == PLAYER_B {
		return chessBoard.BLACK
	}
	return chessBoard.WHITE
}

func GetGameTerminationStatus(board *chessBoard.Board) (bool, bool, bool, string) {
	isDraw := board.IsStalemate()
	isCheckmate := board.IsCheckmate()
	hasGameEnded := isCheckmate || isDraw
	var winner string
	if isCheckmate {
		if board.Turn == chessBoard.BLACK {
			winner = "WHITE"
		} else {
			winner = "BLACK"
		}
	}
	return hasGameEnded, isDraw, isCheckmate, winner
}
