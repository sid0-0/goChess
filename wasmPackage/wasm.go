//go:build js && wasm

package main

import (
	"gochess/chessBoard"
	"syscall/js"
)

func interfaceToJsObject(argsObj interface{}) js.Value {
	if argsObj == nil {
		return js.Undefined()
	}
	switch v := argsObj.(type) {
	case js.Value:
		return v
	case map[string]interface{}:
		obj := js.Global().Get("Object").New()
		for k, val := range v {
			obj.Set(k, interfaceToJsObject(val))
		}
		return obj
	default:
		return js.ValueOf(argsObj)
	}
}

func getAllValidMoves(fen string) map[string][]string {
	// Create a board
	board := chessBoard.New()
	allLegalMoves := map[string][]string{}
	// Load the given fen
	err := board.LoadBoard(fen)
	if err != nil {
		return allLegalMoves
	}
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
	return allLegalMoves
}

func main() {
	chessWASM := js.Global().Get("Object").New()
	chessWASM.Set("marco", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf("polo")
	}))

	chessWASM.Set("getAllValidMoves", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			errorObj := js.Global().Get("Error").New("FEN string is required")
			js.Global().Call("throw", errorObj)
			return nil
		}
		fenString := args[0].String()
		return getAllValidMoves(fenString)
	}))

	js.Global().Set("chessWASM", chessWASM)

	eventDetails := interfaceToJsObject(map[string]interface{}{
		"detail": map[string]interface{}{
			"chessWASM": chessWASM,
		},
	})

	event := js.Global().Get("CustomEvent").New("chessWASMReady", eventDetails)
	js.Global().Call("dispatchEvent", event)

	select {}
}
