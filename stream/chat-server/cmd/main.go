package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type ChatMessage struct {
	StreamID string `json:"stream_id"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     int64  `json:"time"`
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.Mutex
	broadcast = make(chan ChatMessage)
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed: ", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	go client.readPump()
	go client.writePump()
}

func (c *Client) readPump() {
	defer c.conn.Close()
	for {
		var msg ChatMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("reading message error: ", err)
			break
		}
		msg.Time = time.Now().Unix()
		broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		msg, ok := <-c.send
		if !ok {
			return
		}
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("sending message error: ", err)
			break
		}
	}
}

func broadcastMessage(pool *pgxpool.Pool) {
	for msg := range broadcast {
		saveMessage(msg, pool)
		jsonMsg, _ := json.Marshal(msg)
		clientsMu.Lock()
		for client := range clients {
			select {
			case client.send <- jsonMsg:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}

func saveMessage(msg ChatMessage, pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), `
		INSERT INTO public.messages (stream_id, username, text, created_at) VALUES ($1, $2, $3, $4)
	`, msg.StreamID, msg.Username, msg.Text, msg.Time)
	if err != nil {
		log.Println("failed to save message: ", err)
	}
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

	go broadcastMessage(pool)
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Chat server started on :8093")
	log.Fatal(http.ListenAndServe(":8093", nil))
}
