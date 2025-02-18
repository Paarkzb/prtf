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

func (r *RepositoryPostgres) SaveChannel(ctx context.Context, channel models.Channel) (uuid.UUID, error) {
	const op = "Repository.postgres.SaveChannel"

	query := `INSERT INTO public.channels(rf_user_id) VALUES ($1) returning id`
	var id uuid.UUID

	err := r.db.QueryRow(ctx, query, channel.UserID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *RepositoryPostgres) GetAllChannels(ctx context.Context) ([]models.Channel, error) {
	const op = "Repository.postgres.GetAllChannels"

	query := `
		SELECT c.user_id, c.live, c.active_stream
		FROM public.channels as c
	`
	channels := make([]models.Channel, 10)

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	for rows.Next() {
		var channel models.Channel
		err = rows.Scan(&channel.UserID, &channel.Live, &channel.ActiveStreamID)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		channels = append(channels, channel)
	}

	return channels, nil
}
