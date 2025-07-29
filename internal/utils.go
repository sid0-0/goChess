package internal

import (
	"encoding/json"
	"errors"
	"gochess/chessBoard"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/davecgh/go-spew/spew"
)

func LoadAllTemplates(location string) (*template.Template, error) {
	var allTemplates *template.Template
	spew.Println("Loading templates from:", location)
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

func ResolveSquareAndMakeMove(board *chessBoard.Board, fromSquareId string, toSquareId string) error {
	fromSquare := ResolveSquare(board, fromSquareId)
	toSquare := ResolveSquare(board, toSquareId)

	if toSquare == nil || fromSquare == nil {
		return errors.New("invalid square")
	}
	err := board.MakeMove(fromSquare, toSquare)
	return err
}
