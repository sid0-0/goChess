package internal

import (
	"gochess/chessBoard"
	"html/template"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRouter(allTemplates *template.Template, board *chessBoard.Board) *chi.Mux {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(30 * time.Second))

	// Load all routes
	loadRoutes(router, allTemplates, board)

	return router
}

func RunServer() {
	// Create a new chess board
	newBoard := chessBoard.New()

	// Load all templates
	allTemplates, err := LoadAllTemplates("templates")
	if err != nil {
		spew.Println("Template parsing failed", err)
		return
	}

	// Initialize the router
	router := setupRouter(allTemplates, newBoard)

	// Start the server
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		spew.Println("Server failed to start:", err)
	} else {
		spew.Println("Started Go Chess Server at port :8080")
	}
}
