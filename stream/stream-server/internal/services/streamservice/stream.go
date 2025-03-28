package streamservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"
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
	GetChannelTokenByChannelId(ctx context.Context, channelID uuid.UUID) (string, error)
	GetChannelRecordings(ctx context.Context, channelID uuid.UUID) ([]models.Recording, error)
}

type StreamProvider interface {
	StartStream(ctx context.Context, channelID uuid.UUID) (uuid.UUID, error)
	EndStream(ctx context.Context, channelID uuid.UUID, recordPath string, duration time.Duration, posterPath string) (uuid.UUID, error)
	GetActiveChannels(ctx context.Context) ([]models.Channel, error)
	GetRecordingById(ctx context.Context, recordingID uuid.UUID) (models.Recording, error)
}

type RedisProvider interface {
	SetChannels(ctx context.Context, channels []models.Channel) error
	GetChannels(ctx context.Context) ([]models.Channel, error)
}

type StreamService struct {
	log             *zap.SugaredLogger
	channelProvider ChannelProvider
	streamProvider  StreamProvider
	redisProvider   RedisProvider
}

func NewStreamService(log *zap.SugaredLogger, channelProvider ChannelProvider, streamProvider StreamProvider, redisProvider RedisProvider) *StreamService {
	return &StreamService{
		log:             log,
		channelProvider: channelProvider,
		streamProvider:  streamProvider,
		redisProvider:   redisProvider,
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
	const op = "StreamService.StartStream"

	log := s.log.With("op", op, "channelID", channelID)

	streamID, err := s.streamProvider.StartStream(ctx, channelID)
	if err != nil {
		log.Infow("failed to start stream", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	return streamID, nil
}

func (s *StreamService) EndStream(ctx context.Context, channel models.Channel) (uuid.UUID, error) {
	const op = "StreamService.EndStream"

	log := s.log.With("op", op, "channelID", channel.ID)

	cmd := exec.Command("sh", "/var/scripts/save_record.sh", channel.ChannelName)
	recordDirBytes, err := cmd.Output()
	if err != nil {
		log.Infow("failed to save record", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}
	recordDir := string(recordDirBytes)
	recordPath := fmt.Sprintf("%s%s%s%s", recordDir, "/", channel.ChannelName, ".m3u8")

	cmd = exec.Command("sh", "/var/scripts/get_duration.sh", recordDir, channel.ChannelName)
	durationOut, err := cmd.Output()
	if err != nil {
		log.Infow("failed to get duration", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	seconds, err := strconv.ParseFloat(string(durationOut), 64)
	if err != nil {
		log.Infow("failed to parse duration", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}
	duration := time.Duration(seconds * float64(time.Second))

	cmd = exec.Command("sh", "/var/scripts/save_poster.sh", recordDir, channel.ChannelName)
	posterPathBytes, err := cmd.Output()
	if err != nil {
		log.Infow("failed to get duration", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}
	posterPath := string(posterPathBytes)

	streamID, err := s.streamProvider.EndStream(ctx, channel.ID, recordPath, duration, posterPath)
	if err != nil {
		log.Infow("failed to end stream", "err", err)
		return uuid.Nil, fmt.Errorf("%s, %w", op, err)
	}

	return streamID, nil
}

func (s *StreamService) GetActiveChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "StreamService.GetActiveChannels"

	log := s.log.With("op", op)

	channels, err := s.streamProvider.GetActiveChannels(ctx)
	if err != nil {
		log.Infow("failed to get active channels", "err", err)
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return channels, nil
}

func (s *StreamService) GetChannelRecordings(ctx context.Context, channelID uuid.UUID) ([]models.Recording, error) {
	const op = "StreamService.GetChannelRecordings"

	log := s.log.With("op", op, "channelID", channelID)

	recordings, err := s.channelProvider.GetChannelRecordings(ctx, channelID)
	if err != nil {
		log.Infow("failed to get recordings", "err", err)
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return recordings, nil
}

func (s *StreamService) GetRecordingById(ctx context.Context, recordingID uuid.UUID) (models.Recording, error) {
	const op = "StreamService.GetChannelRecordings"

	log := s.log.With("op", op, "recordingID", recordingID)

	recording, err := s.streamProvider.GetRecordingById(ctx, recordingID)
	if err != nil {
		log.Infow("failed to get recording", "err", err)
		return recording, fmt.Errorf("%s, %w", op, err)
	}

	return recording, nil
}

func (s *StreamService) GetChannelTokenByChannelId(ctx context.Context, channelID uuid.UUID) (string, error) {
	const op = "StreamService.GetChannelTokenByChannelId"

	log := s.log.With("op", op, "channelID", channelID)

	channelToken, err := s.channelProvider.GetChannelTokenByChannelId(ctx, channelID)
	if err != nil {
		log.Infow("failed to get channel token", "err", err)
		return "", fmt.Errorf("%s, %w", op, err)
	}

	return channelToken, nil
}
