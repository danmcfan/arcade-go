//go:build js

package internal

import (
	"fmt"
	"syscall/js"
)

func RefreshOnDisconnect() {
	ws := js.Global().Get("WebSocket").New("ws://localhost:8080/ws")

	ws.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("Connected to server...")
		ws.Call("send", "ping")
		return nil
	}))

	ws.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) any {
		ws.Call("send", "ping")
		return nil
	}))

	ws.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("Disconnected from server, reloading...")
		js.Global().Get("location").Call("reload")
		return nil
	}))
}
