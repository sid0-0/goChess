package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (hub *Hub[T]) IsClientInHub(sessionIdToCheck string) (bool, *Client[T], *Pool[T]) {
	for _, pool := range hub.Pools {
		for _, client := range pool.Clients {
			if sessionIdToCheck == client.ID {
				return true, client, pool
			}
		}
	}
	return false, nil, nil
}

func (hub *Hub[T]) NewPool() *Pool[T] {
	newPoolId := uuid.NewString()
	newPool := &Pool[T]{
		ID:         newPoolId,
		Register:   make(chan *Client[T]),
		Unregister: make(chan *Client[T]),
		Clients:    []*Client[T]{},
		Broadcast:  make(chan []byte),
	}

	go func() {
		for {
			select {
			case newClient := <-newPool.Register:
				log.Println(newClient)
				// Check if the client is already in the hub
				for _, client := range newPool.Clients {
					if client.ID == newClient.ID {
						break
					}
				}
				log.Println("Registering new client:", newClient.ID)
				// Add the new client to the pool
				newPool.Clients = append(newPool.Clients, newClient)
				newClient.PoolID = newPool.ID

			case clientToRemove := <-newPool.Unregister:
				log.Println("Unregistering client:", clientToRemove.ID)
				for idx, client := range newPool.Clients {
					if client.ID == clientToRemove.ID {
						log.Println("Found client to remove:", clientToRemove.ID)
						// Remove the client from the pool
						newPool.Clients = append(newPool.Clients[:idx], newPool.Clients[idx+1:]...)

						// cleanup pool if no clients left
						if len(newPool.Clients) == 0 {
							log.Println("No clients left in pool, removing pool:", newPool.ID)
							delete(hub.Pools, newPool.ID)
							return
						}
					}
				}

			case broadcastMessage := <-newPool.Broadcast:
				for _, client := range newPool.Clients {
					client.Send <- broadcastMessage
				}
			}
		}
	}()
	hub.Pools[newPool.ID] = newPool
	return newPool
}

func NewHub[T any]() *Hub[T] {
	newHub := Hub[T]{
		HubId: uuid.NewString(),
		Pools: map[string]*Pool[T]{},
	}

	return &newHub
}

func NewClient[T any](id string, info *T) *Client[T] {
	newClient := Client[T]{
		ID:   id,
		Info: info,
	}

	return &newClient
}

func (client *Client[T]) StartHandlingMessages(conn *websocket.Conn) {
	client.Conn = conn
	client.Receive = make(chan map[string]any)
	client.Send = make(chan []byte)

	go func() {
		for msg := range client.Send {
			if client.Conn != nil {
				if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("Error sending message:", err)
				}
			}
		}
	}()

	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				return
			}
			mappedMessage := map[string]any{}
			if err := json.Unmarshal(p, &mappedMessage); err != nil {
				log.Println("Error unmarshalling message:", err)
				return
			}
			client.Receive <- mappedMessage
		}
	}()
}
