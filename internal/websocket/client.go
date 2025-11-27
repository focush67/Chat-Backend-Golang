package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

/*
	Client now supports structured messages with a Type and Data field.
	We parse the incoming JSON, switch on Message.Type, and send
	the structured commands to the Hub.
*/

type Client struct {
	Hub  *HUB
	Conn *websocket.Conn
	Send chan []byte

	// Will be set after "identify" message
	UserID string
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
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("readPump error:", err)
			break
		}
		//  Parse outer message wrapper: { "type": "...", "data": { ... } }
		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("Invalid JSON from client", err)
			continue
		}

		switch msg.Type {
		case "identify":
			var payload IdentifyPayload
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				log.Println("Invalid identify payload", err)
				continue
			}
			c.UserID = payload.User
			log.Println("User identified as:", c.UserID)

		case "user_list":
			var payload []IdentifyPayload
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				log.Println("Invalid user_list payload", err)
				continue
			}
			response := struct {
				Type string            `json:"type"`
				Data []IdentifyPayload `json:"data"`
			}{
				Type: "user_list",
				Data: c.Hub.GetConnectedUsers(),
			}

			b, _ := json.Marshal(response)
			c.Send <- b

		case "chat_message":
			var payload ChatMessage
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				log.Println("Invalid chat_message payload", err)
				continue
			}
			log.Println("Got this raw message", rawMsg)
			c.Hub.Broadcast <- rawMsg

		case "private_message":
			var payload PrivateMessage
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				log.Println("Invalid private_nessaage payload", err)
				continue
			}
			c.Hub.Broadcast <- rawMsg

		case "join_room":
			var payload JoinRoomPayload
			if err := json.Unmarshal(msg.Data, &payload); err != nil {
				log.Println("Invalid join_room payload", err)
				continue
			}
			// Room Logic will be implemented here

		default:
			log.Println("Unknown message type", msg.Type)
		}
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
