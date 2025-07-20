package ws

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

func (hub *Hub) NewPool() string {
	id := uuid.NewString()
	hub.Clients[id] = []*Client{}
	return id
}

func (hub *Hub) AddClient(client *Client) {
	client.ID = uuid.NewString()
	hub.Register <- client
}

func NewHub() *Hub {
	newHub := Hub{
		HubId:      uuid.NewString(),
		Clients:    make(map[string][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}

	go func() {
		for {
			select {
			case newClient := <-newHub.Register:
				spew.Println(newClient)
				if newClient.ID == "" {
					newClient.ID = uuid.NewString()
				}
				if newClient.PoolID == "" {
					newClient.PoolID = newHub.NewPool()
				}
				spew.Println("Registering new client:", newClient.ID)
				if _, ok := newHub.Clients[newClient.PoolID]; !ok {
					spew.Println("Pool ID not found, creating new pool:", newClient.PoolID)
					newHub.Clients[newClient.PoolID] = []*Client{}
				}
				// Add the new client to the pool
				newHub.Clients[newClient.PoolID] = append(newHub.Clients[newClient.PoolID], newClient)
				// spew.Println("Current clients in pool:", len(newHub.Clients[newClient.PoolID]))
				// spew.Println("Clients in pool:", spew.Sdump(newHub.Clients[newClient.PoolID]))

			case clientToRemove := <-newHub.Unregister:
				spew.Println("Unregistering client:", clientToRemove.ID)
				if _, ok := newHub.Clients[clientToRemove.PoolID]; !ok {
					spew.Println("Pool ID not found:", clientToRemove.PoolID)
					continue
				}
				// Remove the client from the pool
				for _, clients := range newHub.Clients[clientToRemove.PoolID] {
					if clients.ID == clientToRemove.ID {
						spew.Println("Removing client:", clientToRemove.ID)
						newHub.Clients[clientToRemove.PoolID] = append(newHub.Clients[clientToRemove.PoolID][:0], newHub.Clients[clientToRemove.PoolID][1:]...)
						break
					}
				}
				if len(newHub.Clients[clientToRemove.PoolID]) == 0 {
					spew.Println("No clients left in pool, removing pool:", clientToRemove.PoolID)
					delete(newHub.Clients, clientToRemove.PoolID)
					return
				}
			}
		}
	}()

	return &newHub
}
