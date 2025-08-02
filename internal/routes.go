package internal

import (
	"bytes"
	"encoding/json"
	"gochess/chessBoard"
	"gochess/internal/customMiddleware"
	"gochess/ws"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func loadRoutes(router *chi.Mux, wsHub *ws.Hub[ClientInfoType]) {

	// test route
	router.Get("/marco", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte("Polo!"))
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		clientInfo := clientContextData.WebsocketClient.Info
		currentBoard := clientContextData.Board
		pool := clientContextData.Pool
		w.Header().Set("Content-type", "text/html")

		var err error
		if currentBoard == nil {
			err = templates.ExecuteTemplate(w, "Main", nil)
		} else {
			legalMoves := GetLoadLegalMovesJson(currentBoard)
			jsonLegalMoves, _ := json.Marshal(map[string]any{"loadLegalMoves": legalMoves})
			w.Header().Set("HX-Trigger", string(jsonLegalMoves))

			isCheckmate, winner := GetCheckmateAndWinner(currentBoard)
			boardPlayerColor := GetBoardPlayerColorFromPlayerType(clientInfo.Type)
			templateArgs := map[string]any{
				"board":      currentBoard.GetRepresentationalSquares(boardPlayerColor),
				"legalMoves": legalMoves,
				"boardID":    pool.ID,
			}
			if isCheckmate {
				templateArgs["winner"] = winner
			}
			err = templates.ExecuteTemplate(w, "Main", templateArgs)
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

		clientInfo := clientContextData.WebsocketClient.Info
		clientInfo.Type = PLAYER_W

		poolToBoardMap[pool.ID] = newBoard

		pool.Register <- clientContextData.WebsocketClient

		w.Header().Set("Content-type", "text/html")

		boardPlayerColor := GetBoardPlayerColorFromPlayerType(clientInfo.Type)
		legalMoves := GetLoadLegalMovesJson(newBoard)
		templateArgs := map[string]any{
			"board":   newBoard.GetRepresentationalSquares(boardPlayerColor),
			"boardID": pool.ID,
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

	router.Post("/join_game", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		poolToBoardMap := ctx.Value(customMiddleware.PoolToBoardMapContextKey).(customMiddleware.PoolToBoardMap)
		websocketClient := clientContextData.WebsocketClient
		clientInfo := clientContextData.WebsocketClient.Info

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		gameIdToJoin := r.FormValue("gameID")

		var pool *ws.Pool[ClientInfoType]
		for _, poolToCheck := range wsHub.Pools {
			if poolToCheck.ID == gameIdToJoin {
				pool = poolToCheck
			}
		}

		board := poolToBoardMap[gameIdToJoin]

		if board == nil || pool == nil {
			http.Error(w, "Pool not found", http.StatusBadRequest)
			return
		}

		if len(pool.Clients) > 1 {
			clientInfo.Type = SPECTATOR
		} else {
			clientInfo.Type = PLAYER_B
		}

		pool.Register <- websocketClient

		w.Header().Set("Content-type", "text/html")

		boardPlayerColor := GetBoardPlayerColorFromPlayerType(clientInfo.Type)
		legalMoves := GetLoadLegalMovesJson(board)
		templateArgs := map[string]any{
			"board":   board.GetRepresentationalSquares(boardPlayerColor),
			"boardID": pool.ID,
		}
		jsonLegalMoves, _ := json.Marshal(map[string]any{"loadLegalMoves": legalMoves})
		w.Header().Set("HX-Trigger", string(jsonLegalMoves))
		err = templates.ExecuteTemplate(w, "BoardContainer", templateArgs)
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
		clientInfo := clientContextData.WebsocketClient.Info
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

		err = ResolveSquareAndMakeMove(currentBoard, clientInfo.Type, fromSquareId, toSquareId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		legalMoves := GetLoadLegalMovesJson(currentBoard)
		jsonLegalMoves, _ := json.Marshal(map[string]any{"loadLegalMoves": legalMoves})
		w.Header().Set("HX-Trigger", string(jsonLegalMoves))
		boardPlayerColor := GetBoardPlayerColorFromPlayerType(clientInfo.Type)
		templateArgs := map[string]any{
			"board": currentBoard.GetRepresentationalSquares(boardPlayerColor),
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
		log.Println("WebSocket connection requested")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		ctx := r.Context()
		templates := ctx.Value(templatesContextKey).(*template.Template)
		clientContextData := ctx.Value(clientContextDataKey).(*ClientContextData)
		clientContextData.WebsocketClient.StartHandlingMessages(conn)
		clientInfo := clientContextData.WebsocketClient.Info
		pool := clientContextData.Pool
		board := clientContextData.Board

		log.Println("WebSocket connected")

		go func() {
			for msg := range clientContextData.WebsocketClient.Receive {
				if msg["type"] == "move" {
					err = ResolveSquareAndMakeMove(board, clientInfo.Type, msg["from"].(string), msg["to"].(string))
					if err != nil {
						log.Println("Error making move:", err)
						continue
					}

					isCheckmate, winner := GetCheckmateAndWinner(board)
					for _, client := range pool.Clients {
						var buffer bytes.Buffer
						boardPlayerColor := GetBoardPlayerColorFromPlayerType(client.Info.Type)
						templateArgs := map[string]any{
							"board": board.GetRepresentationalSquares(boardPlayerColor),
						}
						if isCheckmate {
							templateArgs["winner"] = winner
						}
						templates.ExecuteTemplate(&buffer, "Board", templateArgs)
						client.Send <- buffer.Bytes()
					}
					legalMoves := GetLoadLegalMovesJson(board)
					pool.Broadcast <- []byte(`{"type": "loadLegalMoves", "data": ` + legalMoves + `}`)
				}
			}
		}()
	})

}
