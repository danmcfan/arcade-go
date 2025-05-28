//go:build js

package internal

import (
	"syscall/js"
)

func RefreshOnDisconnect() {
	ws := js.Global().Get("WebSocket").New("ws://localhost:8080/ws")

	ws.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) any {
		js.Global().Get("location").Call("reload")
		return nil
	}))
}
