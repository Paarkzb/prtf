package streamservice

import (
	"context"
	"videostream/internal/domain/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChannelProvider interface {
	SaveChannel(ctx context.Context, channel models.Channel) (uuid.UUID, error)
	GetAllChannels(ctx context.Context) ([]models.Channel, error)
}

type StreamProvider interface {
}

type StreamService struct {
	log             *zap.SugaredLogger
	channelProvider ChannelProvider
	streamProvider  StreamProvider
}

func NewStreamService(log *zap.SugaredLogger, channelProvider ChannelProvider) *StreamService {
	return &StreamService{
		log:             log,
		channelProvider: channelProvider,
	}
}

func (s *StreamService) SaveChannel(ctx context.Context, channel models.Channel) (uuid.UUID, error) {
	const op = "StreamService.SaveChannel"

	s.log.With("op", op, "userID", channel.UserID)

	var channelID uuid.UUID
	channelID, err := s.channelProvider.SaveChannel(ctx, channel)
	if err != nil {
		s.log.Infow("failed to save channel ", err)
		return uuid.Nil, err
	}

	return channelID, nil
}

func (s *StreamService) GetAllChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "StreamService.GetAllChannels"

	s.log.With("op", op)

	channels, err := s.channelProvider.GetAllChannels(ctx)
	if err != nil {
		s.log.Infow("failed to save channel ", err)
		return nil, err
	}

	return channels, nil
}
