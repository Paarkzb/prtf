package websocketservice

import (
	"chat-server/internal/domain/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessageProvider interface {
	SaveMessage(ctx context.Context, msg models.Message) (uuid.UUID, error)
	// getMessage
}

type RedisProvider interface {
}

type WebsocketService struct {
	log             *zap.SugaredLogger
	Pool            *pool
	messageProvider MessageProvider
	redisProvider   RedisProvider
}

func NewWebsocketService(log *zap.SugaredLogger, messageProvider MessageProvider, redisProvider RedisProvider) *WebsocketService {
	return &WebsocketService{
		log:             log,
		Pool:            newPool(),
		messageProvider: messageProvider,
		redisProvider:   redisProvider,
	}
}

func (s *WebsocketService) Start() {
	for {
		select {
		case client := <-s.Pool.Register:
			if s.Pool.Channels[client.StreamChannelID] == nil {
				s.Pool.Channels[client.StreamChannelID] = make(map[*Client]bool)
			}
			s.Pool.Channels[client.StreamChannelID][client] = true
			s.log.Infow("register", "channels", fmt.Sprint(s.Pool.Channels))
			// client.Poll.Broadcast <- models.Message{StreamID: uuid.Nil, Text: fmt.Sprintf("%s подключился к чату", client.Channel.ChannelName), Time: time.Now(), Channel: client.Channel}

		case client := <-s.Pool.Unregister:
			delete(s.Pool.Channels[client.StreamChannelID], client)
			s.log.Infow("unregister", "channels", fmt.Sprint(s.Pool.Channels))
			// s.log.Println("Size of Connection Pool: ", len(s.Pool.Clients))
			// client.Poll.Broadcast <- models.Message{StreamID: uuid.Nil, Text: "Пользователь отключился", Time: time.Now(), Channel: client.Channel}

		case message := <-s.Pool.Broadcast:
			s.log.Infow("broadcast", "channels", fmt.Sprint(s.Pool.Channels))

			_, err := s.messageProvider.SaveMessage(context.Background(), message)
			if err != nil {
				s.log.Errorw("failed to save message", "err", err)
				return
			}

			s.log.Infow("Send message to all clients of channel", "stream_channel_id", message.StreamChannelId, "text", message.Text)
			for client := range s.Pool.Channels[message.StreamChannelId] {
				if err := client.Conn.WriteJSON(message); err != nil {
					s.log.Errorw("failed to broadcast message", "err", err)
					return
				}
			}
		}
	}
}
