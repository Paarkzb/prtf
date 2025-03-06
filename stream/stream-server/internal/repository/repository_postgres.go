package repository

import (
	"context"
	"fmt"
	"time"
	"videostream/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewRepositoryPostgres(db *pgxpool.Pool) *RepositoryPostgres {
	return &RepositoryPostgres{
		db: db,
	}
}

func (r *RepositoryPostgres) SaveChannel(ctx context.Context, channel models.Channel, streamToken string) (uuid.UUID, error) {
	const op = "Repository.postgres.SaveChannel"

	query := `INSERT INTO public.channels(rf_user_id, channel_token) VALUES ($1, $2) returning id`
	var id uuid.UUID

	err := r.db.QueryRow(ctx, query, channel.RfUserID, streamToken).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *RepositoryPostgres) GetAllChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "Repository.postgres.GetAllChannels"

	query := `
		SELECT c.id, c.rf_user_id, c.channel_name, c.live, c.rf_active_stream_id, c.created_at, c.updated_at
		FROM public.channels as c
	`
	channels := make([]models.Channel, 0)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.RfUserID, &channel.ChannelName, &channel.Live, &channel.RfActiveStreamID, &channel.CreatedAt, &channel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		channels = append(channels, channel)
	}
	rows.Close()

	return channels, nil
}

func (r *RepositoryPostgres) GetChannelById(ctx context.Context, channelID uuid.UUID) (models.Channel, error) {
	const op = "Repository.postgres.GetChannelById"

	query := `
		SELECT c.id, c.rf_user_id, c.channel_name, c.live, c.rf_active_stream_id, c.created_at, c.updated_at
		FROM public.channels as c
		WHERE c.id = $1
	`
	var channel models.Channel

	err := r.db.QueryRow(ctx, query, channelID).Scan(&channel.ID, &channel.RfUserID, &channel.ChannelName, &channel.Live, &channel.RfActiveStreamID, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		return channel, fmt.Errorf("%s:%w", op, err)
	}

	return channel, nil
}

func (r *RepositoryPostgres) GetChannelByUserId(ctx context.Context, userID uuid.UUID) (models.Channel, error) {
	const op = "Repository.postgres.GetChannelByUserId"

	query := `
		SELECT c.id, c.rf_user_id, c.channel_name, c.live, c.rf_active_stream_id, c.channel_token, c.created_at, c.updated_at
		FROM public.channels as c
		WHERE c.rf_user_id = $1
	`
	var channel models.Channel

	err := r.db.QueryRow(ctx, query, userID).Scan(&channel.ID, &channel.RfUserID, &channel.ChannelName, &channel.Live, &channel.RfActiveStreamID, &channel.ChannelToken, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		return channel, fmt.Errorf("%s:%w", op, err)
	}

	return channel, nil
}

func (r *RepositoryPostgres) GetChannelByChannelToken(ctx context.Context, channelToken string) (models.Channel, error) {
	const op = "Repository.postgres.getChannelByChannelToken"

	query := `
		SELECT c.id, c.rf_user_id, c.channel_name, c.live, c.rf_active_stream_id, c.channel_token, c.created_at, c.updated_at
		FROM public.channels as c
		WHERE c.channel_token = $1
	`
	var channel models.Channel

	err := r.db.QueryRow(ctx, query, channelToken).Scan(&channel.ID, &channel.RfUserID, &channel.ChannelName, &channel.Live, &channel.RfActiveStreamID, &channel.ChannelToken, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		return channel, fmt.Errorf("%s:%w", op, err)
	}

	return channel, nil
}

func (r *RepositoryPostgres) GetChannelTokenByChannelId(ctx context.Context, channelID uuid.UUID) (string, error) {
	const op = "Repository.postgres.GetChannelTokenByChannelId"

	query := `
		SELECT c.channel_token
		FROM public.channels as c
		WHERE c.id = $1
	`

	var channelToken string

	err := r.db.QueryRow(ctx, query, channelID).Scan(&channelToken)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return channelToken, nil
}

func (r *RepositoryPostgres) GetActiveChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "Repository.postgres.GetActiveChannels"

	query := `
		SELECT c.id, c.rf_user_id, c.channel_name, c.live, c.rf_active_stream_id, c.icon, c.created_at, c.updated_at
		FROM public.channels as c
		WHERE c.live = true
	`

	var channels []models.Channel

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.RfUserID, &channel.ChannelName, &channel.Live, &channel.RfActiveStreamID, &channel.Icon, &channel.CreatedAt, &channel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		channels = append(channels, channel)
	}
	rows.Close()

	return channels, nil
}

func (r *RepositoryPostgres) GetChannelRecordings(ctx context.Context, channelID uuid.UUID) ([]models.Recording, error) {
	const op = "Repository.postgres.GetChannelRecordings"

	query := `
		SELECT s.id, c.channel_name, s.recording_path, s.created_at, s.duration, s.poster
		FROM public.streams as s
		INNER JOIN public.channels as c ON c.id = s.rf_channel_id AND c.rf_active_stream_id != s.id
		WHERE c.id = $1
	`

	var recordings []models.Recording

	rows, err := r.db.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var recording models.Recording
		err = rows.Scan(&recording.ID, &recording.ChannelName, &recording.Path, &recording.Date, &recording.Duration, &recording.Poster)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		recordings = append(recordings, recording)
	}
	rows.Close()

	return recordings, nil
}

func (r *RepositoryPostgres) GetRecordingById(ctx context.Context, recordingID uuid.UUID) (models.Recording, error) {
	const op = "Repository.postgres.GetRecordingById"

	query := `
		SELECT s.id, s.rf_channel_id, c.channel_name, s.recording_path, s.created_at, s.duration, s.poster
		FROM public.streams as s
		INNER JOIN public.channels as c ON c.id = s.rf_channel_id AND c.rf_active_stream_id != s.id
		WHERE s.id = $1
	`

	var recording models.Recording

	err := r.db.QueryRow(ctx, query, recordingID).Scan(&recording.ID, &recording.ChannelId, &recording.ChannelName, &recording.Path, &recording.Date, &recording.Duration, &recording.Poster)
	if err != nil {
		return recording, fmt.Errorf("%s:%w", op, err)
	}

	return recording, nil
}

func (r *RepositoryPostgres) StartStream(ctx context.Context, channelID uuid.UUID) (uuid.UUID, error) {
	const op = "Repository.postgres.StartStream"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `
		INSERT INTO public.streams(rf_channel_id)
		SELECT c.id
		FROM public.channels as c
		WHERE c.id = $1
		RETURNING public.streams.id
	`
	var streamID uuid.UUID

	err = tx.QueryRow(ctx, query, channelID).Scan(&streamID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	query = `
		UPDATE public.channels SET rf_active_stream_id = $1, live=true WHERE public.channels.id = $2
	`

	_, err = tx.Exec(ctx, query, streamID, channelID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("%s:%w", op, err)
	}

	return streamID, nil
}

func (r *RepositoryPostgres) EndStream(ctx context.Context, channelID uuid.UUID, recordPath string, duration time.Duration, posterPath string) (uuid.UUID, error) {
	const op = "Repository.postgres.EndStream"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `
		UPDATE public.streams
		SET recording_path = $2,
			duration = $3,
			poster = $4
		WHERE public.streams.id = (
			SELECT c.rf_active_stream_id
			FROM public.channels as c
			WHERE c.id = $1
		)
		RETURNING public.streams.id
		
	`
	var streamID uuid.UUID

	err = tx.QueryRow(ctx, query, channelID, recordPath, duration, posterPath).Scan(&streamID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	query = `
		UPDATE public.channels SET rf_active_stream_id = uuid_nil(), live=false WHERE public.channels.id = $1
	`

	_, err = tx.Exec(ctx, query, channelID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("%s:%w", op, err)
	}

	return streamID, nil
}
