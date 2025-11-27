package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

/*
	Upgrader config is responsible for converting an HTTP request into a WebSocket connection.
	CheckOrigin allows all origins here (only for development). Later on we can restrict it for security
*/

var Upgrader = websocket.Upgrader{
	// CheckOrigin prevents cross-site WebSocket hijacking. For now, returning true (allow all endpoints to access it)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
