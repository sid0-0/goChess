package ws

import "github.com/gorilla/websocket"

type Client struct {
	ID         string          // Unique identifier for the client
	Conn       *websocket.Conn // WebSocket connection
	Send       chan []byte     // Channel to send messages to the client
	ClientType string          // Group identifier for the client
	PoolID     string          // Identifier for the board the client is connected to
}

type Hub struct {
	HubId      string               // Unique identifier for the hub
	Clients    map[string][]*Client // Map of connected clients
	Register   chan *Client         // Channel to register new clients
	Unregister chan *Client         // Channel to unregister clients
	Broadcast  chan []byte          // Channel to broadcast messages to all clients
}
