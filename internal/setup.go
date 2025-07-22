package internal

import (
	"context"
	"gochess/internal/customMiddleware"
	"gochess/ws"
	"html/template"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRouter(allTemplates *template.Template, wsHub *ws.Hub) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(30 * time.Second))

	router.Use(customMiddleware.CookieHandler)
	router.Use(customMiddleware.PoolToBoardMapMiddleware)

	// load context in requests
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a new context with the templates and board hubs
			ctx := r.Context()

			poolToBoardMap := r.Context().Value(customMiddleware.PoolToBoardMapContextKey).(customMiddleware.PoolToBoardMap)

			sessionId, _ := r.Cookie(customMiddleware.CookieKey)

			ok, client, clientPool := wsHub.IsClientInHub(sessionId.Value)

			if ok {
				ctx = context.WithValue(ctx, clientContextDataKey, &ClientContextData{
					Board:           poolToBoardMap[clientPool.ID],
					WebsocketClient: client,
					Pool:            clientPool,
				})
			} else {
				ctx = context.WithValue(ctx, clientContextDataKey, &ClientContextData{
					Board:           nil,
					WebsocketClient: ws.NewClient(sessionId.Value),
					Pool:            nil,
				})

			}

			ctx = context.WithValue(ctx, templatesContextKey, allTemplates)
			// Call the next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Load all routes
	loadRoutes(router, wsHub)

	return router
}

func RunServer() {

	// Load all templates
	allTemplates, err := LoadAllTemplates("templates")
	if err != nil {
		spew.Println("Template parsing failed", err)
		return
	}

	// Create a new chess board for local dev
	wsHub := ws.NewHub()

	// Initialize the router
	router := setupRouter(allTemplates, wsHub)

	// Start the server
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		spew.Println("Server failed to start:", err)
	} else {
		spew.Println("Started Go Chess Server at port :8080")
	}
}
