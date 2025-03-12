package handler

import (
	"chat-server/internal/domain/models"
	"chat-server/internal/lib/websocket"
	"chat-server/internal/services/websocketservice"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	log              *zap.SugaredLogger
	websocketService *websocketservice.WebsocketService
}

func NewHandler(log *zap.SugaredLogger, websocketService *websocketservice.WebsocketService) *Handler {
	return &Handler{
		log:              log,
		websocketService: websocketService,
	}
}

func (h *Handler) serverWS(w http.ResponseWriter, r *http.Request) {
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

	client := &websocketservice.Client{
		Channel:         msg.Channel,
		Poll:            h.websocketService.Pool,
		Conn:            conn,
		StreamChannelID: msg.StreamChannelId,
	}

	h.websocketService.Pool.Register <- client
	client.Read()
}

func (h *Handler) InitRoutes() *http.ServeMux {

	go h.websocketService.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", h.serverWS)

	return mux
}
