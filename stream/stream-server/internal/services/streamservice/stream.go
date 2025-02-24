package streamservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"videostream/internal/domain/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChannelProvider interface {
	SaveChannel(ctx context.Context, channel models.Channel, streamToken string) (uuid.UUID, error)
	GetAllChannels(ctx context.Context) ([]models.Channel, error)
	GetChannelById(ctx context.Context, channelID uuid.UUID) (models.Channel, error)
	GetChannelByUserId(ctx context.Context, userId uuid.UUID) (models.Channel, error)
	GetChannelByChannelToken(ctx context.Context, channelToken string) (models.Channel, error)
	GetChannelTokenById(ctx context.Context, channelID uuid.UUID) (string, error)
}

type StreamProvider interface {
	StartStream(ctx context.Context, channelID uuid.UUID) (uuid.UUID, error)
	GetActiveChannels(ctx context.Context) ([]models.Channel, error)
}

type StreamService struct {
	log             *zap.SugaredLogger
	channelProvider ChannelProvider
	streamProvider  StreamProvider
}

func NewStreamService(log *zap.SugaredLogger, channelProvider ChannelProvider, streamProvider StreamProvider) *StreamService {
	return &StreamService{
		log:             log,
		channelProvider: channelProvider,
		streamProvider:  streamProvider,
	}
}

func generateChannelToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return "live_" + hex.EncodeToString(b)
}

func (s *StreamService) SaveChannel(ctx context.Context, channel models.Channel) (uuid.UUID, error) {
	const op = "StreamService.SaveChannel"

	log := s.log.With("op", op, "userID", channel.RfUserID)

	channelToken := generateChannelToken()

	if len(channelToken) == 0 {
		log.Infow("failed to save channel")
		return uuid.Nil, fmt.Errorf("%s, %w", op, errors.New("failed to generate token"))
	}

	var channelID uuid.UUID
	channelID, err := s.channelProvider.SaveChannel(ctx, channel, channelToken)
	if err != nil {
		log.Infow("failed to save channel", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	return channelID, nil
}

func (s *StreamService) GetAllChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "StreamService.GetAllChannels"

	log := s.log.With("op", op)

	channels, err := s.channelProvider.GetAllChannels(ctx)
	if err != nil {
		log.Infow("failed to get channels", "err", err)
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return channels, nil
}

func (s *StreamService) GetChannelById(ctx context.Context, channelID uuid.UUID) (models.Channel, error) {
	const op = "StreamService.GetChannelById"

	log := s.log.With("op", op, "channelID", channelID)

	channel, err := s.channelProvider.GetChannelById(ctx, channelID)
	if err != nil {
		log.Infow("failed to get channel data", "err", err)
		return channel, fmt.Errorf("%s, %w", op, err)
	}

	return channel, nil
}

func (s *StreamService) GetChannelByUserId(ctx context.Context, userID uuid.UUID) (models.Channel, error) {
	const op = "StreamService.GetChannelByUserId"

	log := s.log.With("op", op, "userID", userID)

	channel, err := s.channelProvider.GetChannelByUserId(ctx, userID)
	if err != nil {
		log.Infow("failed to get channel data", "err", err)
		return channel, fmt.Errorf("%s, %w", op, err)
	}

	return channel, nil
}

// func (s *StreamService) GenerateStreamToken(ctx context.Context, channel models.Channel) (string, error) {
// 	const op = "StreamService.GenerateStreamToken"

// 	log := s.log.With("op", op, "userID", channel.RfUserID)

// 	streamToken, err := jwt.NewStreamToken(channel, 1*time.Hour)
// 	if err != nil {
// 		log.Infow("failed to generate stream token", "err", err)
// 		return "", fmt.Errorf("%s, %w", op, err)
// 	}

// 	return streamToken, nil
// }

func (s *StreamService) ValidateStreamToken(ctx context.Context, streamKey string) (models.Channel, error) {
	const op = "StreamService.ValidateStreamToken"

	log := s.log.With("op", op, "streamKey", streamKey)

	channel, err := s.channelProvider.GetChannelByChannelToken(ctx, streamKey)
	if err != nil {
		log.Infow("uncorrect stream token", "err", err)
		return channel, fmt.Errorf("%s, %w", op, err)
	}

	return channel, nil
}

func (s *StreamService) StartStream(ctx context.Context, channelID uuid.UUID) (uuid.UUID, error) {
	const op = "StreamService.ValidateStreamToken"

	log := s.log.With("op", op, "channelID", channelID)

	streamID, err := s.streamProvider.StartStream(ctx, channelID)
	if err != nil {
		log.Infow("failed to start stream", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	return streamID, nil
}

func (s *StreamService) GetActiveChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "StreamService.GetActiveChannels"

	log := s.log.With("op", op)

	channels, err := s.streamProvider.GetActiveChannels(ctx)
	if err != nil {
		log.Infow("failed to start stream", "err", err)
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return channels, nil
}
