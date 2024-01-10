package main

import (
	"fmt"
	"gochess/chessBoard"
	"io/fs"
	"net/http"
	"path/filepath"
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
	return filepath.WalkDir(location, func(path string, _ fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err = allTemplates.New("").Funcs(sprig.FuncMap()).ParseFiles(path)
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
		fmt.Println(len(newBoard.Squares))
		// spew.Dump(newBoard.Squares)
		allTemplates.ExecuteTemplate(w, "Main", map[string]any{"board": newBoard.GetRepresentationalSquares()})
	})

	http.ListenAndServe("127.0.0.1:8080", r)
	fmt.Println("Listening on port 8080")
}
