package internal

import (
	"fmt"
	"gochess/chessBoard"
	"gochess/internal/customMiddleware"
	"gochess/ws"
	"html/template"
	"net/http"
	"slices"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func loadRoutes(router *chi.Mux, wsHub *ws.Hub) {

	// test route
	router.Get("/marco", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte("Polo!"))
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		currentBoard := clientContextData.Board
		w.Header().Set("Content-type", "text/html")

		var err error
		if currentBoard == nil {
			err = templates.ExecuteTemplate(w, "Main", nil)
		} else {
			err = templates.ExecuteTemplate(w, "Main", map[string]any{"board": currentBoard.GetRepresentationalSquares()})
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	router.Get("/start_new_game", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		poolToBoardMap := ctx.Value(customMiddleware.PoolToBoardMapContextKey).(customMiddleware.PoolToBoardMap)

		pool := wsHub.NewPool()
		newBoard := chessBoard.New()

		poolToBoardMap[pool.ID] = newBoard

		pool.Register <- clientContextData.WebSocketData

		w.Header().Set("Content-type", "text/html")

		templateArgs := map[string]any{"board": newBoard.GetRepresentationalSquares()}
		err := templates.ExecuteTemplate(w, "BoardContainer", templateArgs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		err = templates.ExecuteTemplate(w, "HomeActions", templateArgs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	var highlighted *chessBoard.Square

	router.Post("/move/{square}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		currentBoard := clientContextData.Board
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

		templates.ExecuteTemplate(w, "Board", map[string]any{"board": currentBoard.GetRepresentationalSquares()})
	})

	router.Post("/highlight/{square}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		currentBoard := clientContextData.Board
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
				templates.ExecuteTemplate(w, "Square", map[string]any{"data": currentBoard.Squares[i][j], "highlight": slices.Contains(square.LegalMoves, &currentBoard.Squares[i][j])})
			}
		}
	})

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow any origin for dev. Lock this down in production!
			return true
		},
	}

	router.Get("/ws/board", func(w http.ResponseWriter, r *http.Request) {
		spew.Println("WebSocket connection requested")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			spew.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		spew.Println("WebSocket connected")

		for {
			time.Sleep(1 * time.Second)
			err := conn.WriteMessage(websocket.TextMessage, []byte("<div>Ping from server</div>"))
			if err != nil {
				spew.Println("Write error:", err)
				break
			}
		}
	})

}
