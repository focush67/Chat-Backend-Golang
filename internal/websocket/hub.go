package websocket

/*

HUB is the central message router. It keeps track of all connected Clients.
It receives:
	1. Register requests (new Clients)
	2. Unregister requests (disconnected Clients)
	3. Broadcast requests (messages from Clients)
*/

type HUB struct {
	Register   chan *Client     // When a new client connects, its pointer is sent to this channel
	Unregister chan *Client     // When a client disconnects, its pointer goes here
	Broadcast  chan []byte      // Any message that should be Broadcast to all Clients goes here
	Clients    map[*Client]bool /* A map of all active Clients.
	Key = *Client
	Value = bool (true means active)
	*/
}

// NewHUB initialises a new HUB instance with all channels ready.
func NewHUB() *HUB {
	return &HUB{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*Client]bool),
	}
}

/*
	Run starts an infinite loop that listens on the three HUB channels: Register, Unregister and Broadcast
	This goroutine NEVER stops - it is the perpetual event dispatcher
*/

func (h *HUB) RUN() {
	for {
		select {
		// A new client has connected
		case client := <-h.Register:
			h.Clients[client] = true

		// A client has disconnected
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send) // tell client writePump to stop
			}

		// A message needs to be sent to all connected Clients
		case message := <-h.Broadcast:
			for client := range h.Clients {
				// Try sending the message to client.send channel.
				select {
				case client.Send <- message:
					// message successfully queued for client
				default:
					// If client.send is full or blocked, disconnect it
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}

	}
}
