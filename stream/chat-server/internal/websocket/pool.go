package websocket

import (
	"chat-server/internal/domain/models"
	"log"

	"github.com/google/uuid"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
	Channels   map[uuid.UUID]map[*Client]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client, 1),
		Unregister: make(chan *Client, 1),
		Broadcast:  make(chan models.Message, 1),
		Channels:   make(map[uuid.UUID]map[*Client]bool),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			if p.Channels[client.StreamChannelID] == nil {
				p.Channels[client.StreamChannelID] = make(map[*Client]bool)
			}
			p.Channels[client.StreamChannelID][client] = true
			log.Println("register.channels", p.Channels)
			// client.Poll.Broadcast <- models.Message{StreamID: uuid.Nil, Text: fmt.Sprintf("%s подключился к чату", client.Channel.ChannelName), Time: time.Now(), Channel: client.Channel}

		case client := <-p.Unregister:
			delete(p.Channels[client.StreamChannelID], client)
			log.Println("unregister.channels", p.Channels)
			// log.Println("Size of Connection Pool: ", len(p.Clients))
			// client.Poll.Broadcast <- models.Message{StreamID: uuid.Nil, Text: "Пользователь отключился", Time: time.Now(), Channel: client.Channel}

		case message := <-p.Broadcast:
			log.Println("broadcast.channels", p.Channels)
			// save message to db
			log.Println("Send message to all clients of channel", message.StreamChannelId, message.Text)
			for client := range p.Channels[message.StreamChannelId] {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
