package game

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Pool   *Pool
	Player *Player
}

func NewClient(conn *websocket.Conn, pool *Pool, player *Player) *Client {
	return &Client{
		ID:     uuid.New().String(),
		Conn:   conn,
		Pool:   pool,
		Player: player,
	}
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		log.Printf("Message Received: %+v\n", message)

		// ticker := time.NewTicker(time.Millisecond * 16)

		// prevUpdate := time.Now()
		// for {
		// 	select {
		// 	case <-ticker.C:
		// 		dt := float64(time.Since(prevUpdate).Milliseconds()) / 1000
		// 		prevUpdate = time.Now()

		// 		player.update(dt)

		// 		game.updateBullets(dt)

		// 		game.checkCollisions(player)

		// 		log.Println("write state in gourutine")
		// 		game.writeMessage("state", player)
		// 	case <-ctx.Done():
		// 		ticker.Stop()
		// 		return
		// 	}
		// }
	}
}
