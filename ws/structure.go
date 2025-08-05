package ws

import "github.com/gorilla/websocket"

type Client[ClientInfoType any] struct {
	ID      string              // Unique identifier for the client
	Conn    *websocket.Conn     // WebSocket connection
	Send    chan []byte         // Channel to send messages to the client
	Receive chan map[string]any // Channel to receive messages from the client
	PoolID  string              // Identifier for the board the client is connected to
	Info    *ClientInfoType     // Group identifier for the client
}

type Pool[ClientInfoType any] struct {
	ID         string                       // Unique identifier for the pool
	Clients    []*Client[ClientInfoType]    // Map of clients in the pool
	Register   chan *Client[ClientInfoType] // Channel to register new clients
	Unregister chan *Client[ClientInfoType] // Channel to unregister clients
	Broadcast  chan []byte                  // Channel to broadcast messages to all clients in the pool
}

type Hub[ClientInfoType any] struct {
	HubId string                           // Unique identifier for the hub
	Pools map[string]*Pool[ClientInfoType] // Map of pools, each containing clients
}
