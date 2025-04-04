package websocketservice

import (
	"chat-server/internal/domain/models"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Channel       models.Channel
	Conn          *websocket.Conn
	Poll          *pool
	StreamChannelID uuid.UUID
}

func (c *Client) Read() {
	defer func() {
		c.Poll.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var message models.Message
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return
		}

		c.Poll.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)
	}
}
