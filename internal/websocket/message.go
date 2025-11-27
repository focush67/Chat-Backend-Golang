package websocket

import "encoding/json"

/*
	Message is the top-level structure for every incoming WebSocket message.
	Every Message MUST have:
		1. Type ---> tells the server what action to perform
		2. Data ---> JSON payload specific to that type

	Example:
	{
		"type":"chat_message",
		"data":{
			"room":"general",
			"content":"Hello World"
		}
	}
*/

type Message struct {
	Type string          `json:"type"` // e.g. "chat_message", "join_room"
	Data json.RawMessage `json:"data"` // raw JSON for flexible decoding later
}

/*
	ChatMessage is used when the client sends a chat message in a room.
*/

type ChatMessage struct {
	Room    string `json:"room"`    // which room the user is sending to
	Content string `json:"content"` // the actual text
}

/*
	PrivateMessage is used when a client wants to send a message to another user directly.
*/

type PrivateMessage struct {
	Receiver string `json:"receiver"` // userID of receiver
	Content  string `json:"content"`  // actual message
}

/*
	JoinRoomPayload is sent by a client to join a room
*/

type JoinRoomPayload struct {
	Room string `json:"room"`
}

/*
	LeaveRoomPayload is sent when leaving a room
*/

type LeaveRoomPayload struct {
	Room string `json:"room"`
}

/*
	TypingPayload is used when user starts typing (will be used later while enhancing)
*/

type TypingPayload struct {
	Room string `json:"room"`
}

/*
	IdentifyPayload is used when client identifies itself on connect:
	{
		"type":"identify",
		"data":{"user":"sparsh"}
	}
*/

type IdentifyPayload struct {
	User string `json:"user"`
}

/*
	ConnectedUsers is used to give out a list of all connected users
	{
		type:"user_list",
		"data":[{"user":"sparsh"},{"user":"priyanshu"}]
	}
*/

type ConnectedUsers struct {
	Data []IdentifyPayload `json:"data"`
}
