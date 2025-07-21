package ws

import "github.com/gorilla/websocket"

type Client struct {
	ID         string          // Unique identifier for the client
	Conn       *websocket.Conn // WebSocket connection
	Send       chan []byte     // Channel to send messages to the client
	ClientType string          // Group identifier for the client
	PoolID     string          // Identifier for the board the client is connected to
}

type Pool struct {
	ID         string       // Unique identifier for the pool
	Clients    []*Client    // Map of clients in the pool
	Register   chan *Client // Channel to register new clients
	Unregister chan *Client // Channel to unregister clients
}

type Hub struct {
	HubId string  // Unique identifier for the hub
	Pools []*Pool // Map of pools, each containing clients
}
