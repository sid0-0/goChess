package main

import (
	"fmt"
	"gochess/chessBoard"
	"io/fs"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var allTemplates *template.Template

func ParseAllTemplates(location string) error {
	allTemplates = template.New("").Funcs(sprig.FuncMap())
	return filepath.WalkDir(location, func(path string, _ fs.DirEntry, _ error) error {
		if strings.HasSuffix(path, ".html") {
			_, err := allTemplates.New("").Funcs(sprig.FuncMap()).ParseFiles(path)
			return err
		}
		return nil
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	err := ParseAllTemplates("templates")
	if err != nil {
		fmt.Println("Template parsing failed", err)
	}

	newBoard := chessBoard.New()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		allTemplates.ExecuteTemplate(w, "Main", map[string]any{"board": newBoard.GetRepresentationalSquares()})
	})

	var highlighted *chessBoard.Square

	r.Post("/move/{square}", func(w http.ResponseWriter, r *http.Request) {
		squareId := chi.URLParam(r, "square")
		fmt.Println(r.Body)
		defer r.Body.Close()
		if len(squareId) != 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ri, ci := int(squareId[1]-'1'), int(squareId[0]-'a')
		square := newBoard.GetSquare(ri, ci)

		if square == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Invalid square"))
			return
		}

		err = newBoard.MakeMove(highlighted, square)

		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte(err.Error()))
		// 	return
		// }

		allTemplates.ExecuteTemplate(w, "Main", map[string]any{"board": newBoard.GetRepresentationalSquares()})
	})

	r.Post("/highlight/{square}", func(w http.ResponseWriter, r *http.Request) {
		squareId := chi.URLParam(r, "square")
		fmt.Println(r.Body)
		defer r.Body.Close()
		if len(squareId) != 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ri, ci := int(squareId[1]-'1'), int(squareId[0]-'a')
		square := newBoard.GetSquare(ri, ci)

		if square == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		highlighted = square
		// hx-include will allow replacing only highlighted and to-be-highlighted squares but I couldn't get it working
		for i := range newBoard.Squares {
			for j := range newBoard.Squares[i] {
				allTemplates.ExecuteTemplate(w, "Square", map[string]any{"data": newBoard.Squares[i][j], "highlight": slices.Contains(square.LegalMoves, &newBoard.Squares[i][j])})
			}
		}
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Server failed to start:", err)
	} else {
		fmt.Println("Started Go Chess Server at port :8080")
	}
}
