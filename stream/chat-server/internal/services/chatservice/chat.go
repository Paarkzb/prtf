package chatservice

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

type ChatService struct {
	log             *zap.SugaredLogger
	messageProvider MessageProvider
	redisProvider   RedisProvider
}

func NewChatService(log *zap.SugaredLogger, messageProvider MessageProvider, redisProvider RedisProvider) *ChatService {
	return &ChatService{
		log:             log,
		messageProvider: messageProvider,
		redisProvider:   redisProvider,
	}
}

func (s *ChatService) SaveMessage(ctx context.Context, msg models.Message) (uuid.UUID, error) {
	const op = "ChatService.SaveMessage"

	log := s.log.With("op", op, "channelID", msg.Channel.ID)

	msgID, err := s.messageProvider.SaveMessage(ctx, msg)
	if err != nil {
		log.Infow("failed to save message", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	return msgID, nil
}
