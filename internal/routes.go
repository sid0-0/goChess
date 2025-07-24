package internal

import (
	"encoding/json"
	"gochess/chessBoard"
	"gochess/internal/customMiddleware"
	"gochess/ws"
	"html/template"
	"net/http"
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
			legalMoves := GetLoadLegalMovesJson(currentBoard)
			jsonLegalMoves, _ := json.Marshal(map[string]any{"loadLegalMoves": legalMoves})
			w.Header().Set("HX-Trigger", string(jsonLegalMoves))
			err = templates.ExecuteTemplate(w, "Main", map[string]any{
				"board":      currentBoard.GetRepresentationalSquares(),
				"legalMoves": legalMoves,
			})
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

		pool.Register <- clientContextData.WebsocketClient

		w.Header().Set("Content-type", "text/html")

		legalMoves := GetLoadLegalMovesJson(newBoard)
		templateArgs := map[string]any{
			"board": newBoard.GetRepresentationalSquares(),
		}
		jsonLegalMoves, _ := json.Marshal(map[string]any{"loadLegalMoves": legalMoves})
		w.Header().Set("HX-Trigger", string(jsonLegalMoves))
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

	router.Post("/move", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		currentBoard := clientContextData.Board

		// get data from request
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		fromSquareId := r.FormValue("from")
		toSquareId := r.FormValue("to")
		defer r.Body.Close()
		if len(toSquareId) != 2 || len(fromSquareId) != 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ri, ci := int(fromSquareId[1]-'1'), int(fromSquareId[0]-'a')
		fromSquare := currentBoard.GetSquare(ri, ci)
		ri, ci = int(toSquareId[1]-'1'), int(toSquareId[0]-'a')
		toSquare := currentBoard.GetSquare(ri, ci)

		if toSquare == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Invalid square"))
			return
		}

		err = currentBoard.MakeMove(fromSquare, toSquare)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		legalMoves := GetLoadLegalMovesJson(currentBoard)
		jsonLegalMoves, _ := json.Marshal(map[string]any{"loadLegalMoves": legalMoves})
		w.Header().Set("HX-Trigger", string(jsonLegalMoves))
		templateArgs := map[string]any{
			"board": currentBoard.GetRepresentationalSquares(),
		}
		templates.ExecuteTemplate(w, "Board", templateArgs)
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

		clientContextData := r.Context().Value(clientContextDataKey).(*ClientContextData)
		clientContextData.WebsocketClient.StartHandlingMessages(conn)

		spew.Println("WebSocket connected")

		go func() {
			for msg := range clientContextData.WebsocketClient.Receive {
				spew.Println("Received message:", string(msg))
			}
		}()

		for {
			clientContextData.WebsocketClient.Send <- []byte("Hello from the server!")
			time.Sleep(1 * time.Second)
		}
	})

}
