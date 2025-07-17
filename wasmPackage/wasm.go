//go:build js && wasm

package chessBoard

import "syscall/js"

func main() {
	js.Global().Set("hello", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		name := args[0].String()
		msg := "Hello, " + name
		js.Global().Get("console").Call("log", msg)
		return js.ValueOf(msg)
	}))
	select {}
}
