package routes

import (
	"chat/internal/handlers"
	"chat/internal/websocket"
	"net/http"
)

// NewRouter initialises all HTTP routes.
func NewRoutes(hub *websocket.HUB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWebsocket(hub, w, r)
	})

	return mux
}
