//go:build js && wasm

package main

import (
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

func main() {
	chessWASM := js.Global().Get("Object").New()
	chessWASM.Set("marco", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf("polo")
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
