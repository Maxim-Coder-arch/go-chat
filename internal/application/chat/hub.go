package chat

import (
	"sync"

	"chat-app/internal/domain"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[Client]struct{}),
	}
}

func (h *Hub) Register(client Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = struct{}{}
}

func (h *Hub) Unregister(client Client) {
	h.mu.Lock()

	_, exists := h.clients[client]
	if exists {
		delete(h.clients, client)
	}

	h.mu.Unlock()

	if exists {
		_ = client.Close()
	}
}

func (h *Hub) Broadcast(message domain.Message) {
	h.mu.RLock()

	failedClients := make([]Client, 0)

	for client := range h.clients {
		if err := client.Send(message); err != nil {
			failedClients = append(failedClients, client)
		}
	}

	h.mu.RUnlock()

	for _, client := range failedClients {
		h.Unregister(client)
	}
}
