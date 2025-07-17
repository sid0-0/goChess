//go:build js && wasm

package main

import (
	"gochess/chessBoard"
	"syscall/js"
)

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

func consoleLog(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}

func main() {

	chessWASM := map[string]interface{}{
		"marco": (func(this js.Value, args []js.Value) interface{} {
			return js.ValueOf("polo")
		}),
		"getAllValidMoves": (func(this js.Value, args []js.Value) interface{} {
			if len(args) < 1 {
				errorObj := js.Global().Get("Error").New("FEN string is required")
				js.Global().Call("throw", errorObj)
				return nil
			}
			fenString := args[0].String()
			return getAllValidMoves(fenString)
		}),
	}

	chessWASMWithPanicWrapper := map[string]interface{}{}
	for key, value := range chessWASM {
		if fn, ok := value.(func(this js.Value, args []js.Value) interface{}); ok {
			chessWASMWithPanicWrapper[key] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				defer func() {
					if r := recover(); r != nil {
						errorObj := js.Global().Get("Error").New("Panic occurred: " + r.(string))
						js.Global().Call("throw", errorObj)
					}
				}()
				return fn(this, args)
			})
		} else {
			chessWASMWithPanicWrapper[key] = value
		}
	}

	js.Global().Set("chessWASM", chessWASMWithPanicWrapper)

	eventDetails := (map[string]interface{}{
		"detail": map[string]interface{}{
			"chessWASM": chessWASMWithPanicWrapper,
		},
	})

	event := js.Global().Get("CustomEvent").New("chessWASMReady", eventDetails)
	js.Global().Call("dispatchEvent", event)

	select {}
}
