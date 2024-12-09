package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"survio/pkg/game"
	"survio/pkg/websocket"
	"time"
)

// serveWS is a HTTP Handler that the has the Game that allows connections
func serveWS(g *game.Game, w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Begin by upgrading the HTTP request
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	random := rand.New(rand.NewSource(time.Now().Unix()))
	player := game.NewPlayer("", &game.Position{
		X: float64(random.Int31n(game.WorldWidth-2*game.MinWidth) + game.MinWidth),
		Y: float64(random.Int31n(game.WorldHeight-2*game.MinHeight) + game.MinHeight),
	}, 0)
	client := game.NewClient(conn, g.Pool, player)
	g.Pool.Register <- client

	client.Read()

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Millisecond * 16)

		prevUpdate := time.Now()
		for {
			select {
			case <-ticker.C:
				dt := float64(time.Since(prevUpdate).Milliseconds()) / 1000
				prevUpdate = time.Now()

				player.update(dt)

				game.UpdateBullets(dt)

				game.checkCollisions(player)

				log.Println("write state in gourutine")
				game.writeMessage("state", player)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}(ctx)

	// handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// handle message
		var event game.Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println(err)
			continue
		}

		handleMessages(event, player)
	}

}

func handleMessages(event game.Event, player *game.Player) {
	switch event.Type {
	case "login":
		player.Name = event.Payload
		g.writeMessage("start", player)

	case "keydown":
		direction := event.Payload
		// log.Println("Key pressed", event.Payload)
		switch direction {
		case "left":
			player.Kyes.A = true
		case "right":
			player.Kyes.D = true
		case "forward":
			player.Kyes.W = true
		case "back":
			player.Kyes.S = true
		case "space":
			player.Kyes.Space = true
		}
	case "keyup":
		direction := event.Payload
		// log.Println("Key unpressed", event.Payload)
		switch direction {
		case "left":
			player.Kyes.A = false
		case "right":
			player.Kyes.D = false
		case "forward":
			player.Kyes.W = false
		case "back":
			player.Kyes.S = false
		case "space":
			player.Kyes.Space = false
		}
	}

	// g.writeState(player)
}

func main() {
	setup()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// setupAPI will start all Routes and their Handlers
func setup() {
	game := game.NewGame()
	go game.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(game, w, r)
	})
}
