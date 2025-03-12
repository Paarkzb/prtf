package websocketservice

import (
	"chat-server/internal/domain/models"

	"github.com/google/uuid"
)

type pool struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
	Channels   map[uuid.UUID]map[*Client]bool
}

func newPool() *pool {
	return &pool{
		Register:   make(chan *Client, 1),
		Unregister: make(chan *Client, 1),
		Broadcast:  make(chan models.Message, 1),
		Channels:   make(map[uuid.UUID]map[*Client]bool),
	}
}
