package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

/*
	Client represents a single WebSocket user connected to the server. It contains:
		1. hub: pointer to the HUB (so it can broadcast messages)
		2. conn: the actual WebSocket connection
		3. send: a channel where outgoing messages are queued
*/

type Client struct {
	Hub  *HUB
	Conn *websocket.Conn
	Send chan []byte
}

/*
	readPump continuously reads message FROM the WebSocket. Everything the user sends goes to:
		client.hub.broadcast <- message

	It ends when the user disconnects OR an error happens
*/

func (c *Client) ReadPump() {
	// Ensure cleanup happens when function exists
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("readPump error:", err)
			break
		}
		// Send received message into HUB broadcast channel
		c.Hub.Broadcast <- msg
	}
}

/*
	writePump continuously listens on client.send channel and writes messages TO the WebSocket. If client.send channel is closed, the connection ends.
*/

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("writePump error:", err)
			break
		}
	}
}
