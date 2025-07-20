package internal

import (
	"fmt"
	"gochess/chessBoard"
	"html/template"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
)

func loadRoutes(router *chi.Mux, allTemplates *template.Template, currentBoard *chessBoard.Board) {

	// test route
	router.Get("/marco", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte("Polo!"))
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")

		err := allTemplates.ExecuteTemplate(w, "Main", map[string]any{"board": currentBoard.GetRepresentationalSquares()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	var highlighted *chessBoard.Square

	router.Post("/move/{square}", func(w http.ResponseWriter, r *http.Request) {
		squareId := chi.URLParam(r, "square")
		fmt.Println(r.Body)
		defer r.Body.Close()
		if len(squareId) != 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ri, ci := int(squareId[1]-'1'), int(squareId[0]-'a')
		square := currentBoard.GetSquare(ri, ci)

		if square == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Invalid square"))
			return
		}

		// err := newBoard.MakeMove(highlighted, square)
		currentBoard.MakeMove(highlighted, square)

		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte(err.Error()))
		// 	return
		// }

		allTemplates.ExecuteTemplate(w, "Main", map[string]any{"board": currentBoard.GetRepresentationalSquares()})
	})

	router.Post("/highlight/{square}", func(w http.ResponseWriter, r *http.Request) {
		squareId := chi.URLParam(r, "square")
		fmt.Println(r.Body)
		defer r.Body.Close()
		if len(squareId) != 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ri, ci := int(squareId[1]-'1'), int(squareId[0]-'a')
		square := currentBoard.GetSquare(ri, ci)

		if square == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		highlighted = square
		// hx-include will allow replacing only highlighted and to-be-highlighted squares but I couldn't get it working
		for i := range currentBoard.Squares {
			for j := range currentBoard.Squares[i] {
				allTemplates.ExecuteTemplate(w, "Square", map[string]any{"data": currentBoard.Squares[i][j], "highlight": slices.Contains(square.LegalMoves, &currentBoard.Squares[i][j])})
			}
		}
	})
}
