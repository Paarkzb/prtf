package websocket

import (
	"chat-server/internal/domain/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
	Clients    map[*Client]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client, 1),
		Unregister: make(chan *Client, 1),
		Broadcast:  make(chan models.Message, 1),
		Clients:    make(map[*Client]bool),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			log.Println("Size of Connection Pool: ", len(p.Clients))
			client.Poll.Broadcast <- models.Message{StreamID: uuid.Nil, Text: fmt.Sprintf("%s подключился к чату", client.Channel.ChannelName), Time: time.Now(), Channel: client.Channel}

		case client := <-p.Unregister:
			delete(p.Clients, client)
			log.Println("Size of Connection Pool: ", len(p.Clients))
			// client.Poll.Broadcast <- models.Message{StreamID: uuid.Nil, Text: "Пользователь отключился", Time: time.Now(), Channel: client.Channel}

		case message := <-p.Broadcast:
			// save message to db
			log.Println("Send message to all clients", message)
			for client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
