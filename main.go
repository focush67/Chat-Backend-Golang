package main

import (
	"chat/internal/routes"
	"chat/internal/websocket"
	"log"
	"net/http"
)

func main() {
	hub := websocket.NewHUB()
	go hub.RUN()

	router := routes.NewRoutes(hub)

	log.Println("WebSocket server running on: 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
