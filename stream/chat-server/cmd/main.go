package main

import (
	"chat-server/internal/domain/models"
	"chat-server/internal/websocket"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func serverWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket Endpoint Hit")

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	var msg models.Message
	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Println(err)
		return
	}

	client := &websocket.Client{
		Channel:         msg.Channel,
		Poll:            pool,
		Conn:            conn,
		StreamChannelID: msg.StreamChannelId,
	}

	pool.Register <- client
	client.Read()
}

// func saveMessage(msg models.Message, pool *pgxpool.Pool) {
// 	_, err := pool.Exec(context.Background(), `
// 		INSERT INTO public.messages (stream_id, username, text, created_at) VALUES ($1, $2, $3, $4)
// 	`, msg.StreamID, msg.Username, msg.Text, msg.Time)
// 	if err != nil {
// 		log.Println("failed to save message: ", err)
// 	}
// }

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWS(pool, w, r)
	})
}

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("connecting to postgres error: ", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("failed to connect postgres: ", err)
	}

	setupRoutes()

	err = http.ListenAndServe(":8093", nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Println("Chat server started on port 8093")
}
