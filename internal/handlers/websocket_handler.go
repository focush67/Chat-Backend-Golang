package handlers

import (
	"chat/internal/websocket"
	"net/http"
)

/*
	ServeWebsocket is the entry point for WebSocket connections. It is a trigger whenever a client hits /ws on the server.
	These are the steps:
		1. Upgrade HTTP request -> WebSocket connection
		2. Create new Client struct
		3. Register client with HUB
		4. Start readPump -> writePump goroutines
*/

func ServeWebsocket(hub *websocket.HUB, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	// Create a new client
	client := &websocket.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// Register client with Hub
	hub.Register <- client

	// Start pumps
	go client.WritePump()
	go client.ReadPump()
}
