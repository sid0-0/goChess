package internal

import (
	"encoding/json"
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
		if strings.HasSuffix(path, ".gohtml") {
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

func GetLoadLegalMovesJson(board *chessBoard.Board) string {
	if board != nil {
		allLegalMoves := map[string][]string{}
		// collect all legal moves in an object
		for _, row := range board.Squares {
			for _, square := range row {
				squareNotation := square.File + square.Rank
				legalMovesForSquare := []string{}
				for _, square := range square.LegalMoves {
					legalMovesForSquare = append(legalMovesForSquare, square.File+square.Rank)
				}
				allLegalMoves[squareNotation] = legalMovesForSquare
			}
		}
		dataMap := map[string]interface{}{
			"loadLegalMoves": allLegalMoves,
		}
		dataMapJson, err := json.Marshal(dataMap)
		if err == nil {
			return string(dataMapJson)
		}
	}
	return ""
}
