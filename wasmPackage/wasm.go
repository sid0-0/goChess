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
	if fen != "" {
		err := board.LoadBoard(fen)
		if err != nil {
			return allLegalMoves
		}
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
func consoleError(args ...interface{}) {
	js.Global().Get("console").Call("error", args...)
}

func main() {

	chessWASM := map[string]interface{}{
		"marco": (func(this js.Value, args []js.Value) interface{} {
			return "polo"
		}),
		"getAllValidMoves": (func(this js.Value, args []js.Value) interface{} {
			returnValue := map[string][]string{}
			if len(args) < 1 {
				returnValue = getAllValidMoves("")
			} else {
				fenString := args[0].String()
				returnValue = getAllValidMoves(fenString)
			}

			// js wasm does not support map[string][]string directly, so we need to convert it
			// to map[string]interface{} that js wasm can understand
			jsCompatibleReturnValue := make(map[string]interface{}, len(returnValue))
			for k, v := range returnValue {
				slice := make([]interface{}, len(v))
				for i, val := range v {
					slice[i] = val
				}
				jsCompatibleReturnValue[k] = slice
			}
			return jsCompatibleReturnValue
		}),
	}

	chessWASMWithPanicWrapper := map[string]interface{}{}
	for key, value := range chessWASM {
		if fn, ok := value.(func(this js.Value, args []js.Value) interface{}); ok {
			chessWASMWithPanicWrapper[key] = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				defer func() {
					if r := recover(); r != nil {
						errorObj := js.Global().Get("Error").New("Panic occurred: " + r.(string))
						consoleError(errorObj)
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
