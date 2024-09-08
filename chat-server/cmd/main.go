package main

import (
	"fmt"
	"log"
	"net/http"

	"chat-server/pkg/websocket"
)

// WS endpoint
func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket Endpoint Hit")
	// upgrade connection
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)

	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
}

func main() {
	log.Println("Chat App v0.01")

	setupRoutes()

	err := http.ListenAndServe(":8071", nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Println("Chat App started on port 8071")
}
