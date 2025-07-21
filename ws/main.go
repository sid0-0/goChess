package ws

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
)

func (hub *Hub) IsClientInHub(sessionIdToCheck string) (bool, *Client, *Pool) {
	for _, pool := range hub.Pools {
		for _, client := range pool.Clients {
			if sessionIdToCheck == client.ID {
				return true, client, pool
			}
		}
	}
	return false, nil, nil
}

func (hub *Hub) NewPool() string {
	newPoolId := uuid.NewString()
	newPool := &Pool{
		ID:         newPoolId,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    []*Client{},
	}

	go func() {
		for {
			select {
			case newClient := <-newPool.Register:
				spew.Println(newClient)
				// Check if the client is already in the hub
				for _, client := range newPool.Clients {
					if client.ID == newClient.ID {
						break
					}
				}
				spew.Println("Registering new client:", newClient.ID)
				// Add the new client to the pool
				newPool.Clients = append(newPool.Clients, newClient)

			case clientToRemove := <-newPool.Unregister:
				spew.Println("Unregistering client:", clientToRemove.ID)
				for idx, client := range newPool.Clients {
					if client.ID == clientToRemove.ID {
						spew.Println("Found client to remove:", clientToRemove.ID)
						// Remove the client from the pool
						newPool.Clients = append(newPool.Clients[:idx], newPool.Clients[idx+1:]...)

						// cleanup pool if no clients left
						if len(newPool.Clients) == 0 {
							spew.Println("No clients left in pool, removing pool:", newPool.ID)
							for idx, hubPool := range hub.Pools {
								if hubPool.ID == newPool.ID {
									spew.Println("Found pool to remove:", hubPool.ID)
									// Remove the pool from the hub
									hub.Pools = append(hub.Pools[:idx], hub.Pools[idx+1:]...)
								}
							}
							return
						}
					}
				}
			}
		}
	}()
	hub.Pools = append(hub.Pools, newPool)
	return newPoolId
}

func (pool *Pool) AddClient(client *Client) {
	client.ID = uuid.NewString()
	pool.Register <- client
}

func NewHub() *Hub {
	newHub := Hub{
		HubId: uuid.NewString(),
		Pools: []*Pool{},
	}

	return &newHub
}
