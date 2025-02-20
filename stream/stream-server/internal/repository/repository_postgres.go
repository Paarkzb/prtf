package repository

import (
	"context"
	"fmt"
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
		SELECT c.id, c.rf_user_id, c.live, c.rf_active_stream_id, c.created_at, c.updated_at
		FROM public.channels as c
	`
	channels := make([]models.Channel, 0)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.ID, &channel.RfUserID, &channel.Live, &channel.RfActiveStreamID, &channel.CreatedAt, &channel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

func (r *RepositoryPostgres) GetChannelById(ctx context.Context, channelID uuid.UUID) (models.Channel, error) {
	const op = "Repository.postgres.GetChannelById"

	query := `
		SELECT c.id, c.rf_user_id, c.live, c.rf_active_stream_id, c.created_at, c.updated_at
		FROM public.channels as c
		WHERE c.id = $1
	`
	var channel models.Channel

	err := r.db.QueryRow(ctx, query, channelID).Scan(&channel.ID, &channel.RfUserID, &channel.Live, &channel.RfActiveStreamID, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		return channel, fmt.Errorf("%s:%w", op, err)
	}

	return channel, nil
}

func (r *RepositoryPostgres) GetChannelByUserId(ctx context.Context, userID uuid.UUID) (models.Channel, error) {
	const op = "Repository.postgres.GetChannelByUserId"

	query := `
		SELECT c.id, c.rf_user_id, c.live, c.rf_active_stream_id, c.channel_token, c.created_at, c.updated_at
		FROM public.channels as c
		WHERE c.rf_user_id = $1
	`
	var channel models.Channel

	err := r.db.QueryRow(ctx, query, userID).Scan(&channel.ID, &channel.RfUserID, &channel.Live, &channel.RfActiveStreamID, &channel.ChannelToken, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		return channel, fmt.Errorf("%s:%w", op, err)
	}

	return channel, nil
}

func (r *RepositoryPostgres) GetChannelTokenById(ctx context.Context, channelID uuid.UUID) (string, error) {
	const op = "Repository.postgres.GetChannelTokenById"

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
