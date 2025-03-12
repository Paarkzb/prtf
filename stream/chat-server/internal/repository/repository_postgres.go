package repository

import (
	"chat-server/internal/domain/models"
	"context"
	"fmt"

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

func (r *RepositoryPostgres) SaveMessage(ctx context.Context, msg models.Message) (uuid.UUID, error) {
	const op = "Repository.postgres.SaveMessage"

	query := `INSERT INTO public.messages (rf_stream_channel_id, rf_stream_id, rf_channel_id, text, created_at) VALUES ($1, $2, $3, $4, $5) returning id`
	var id uuid.UUID

	err := r.db.QueryRow(ctx, query, msg.StreamChannelId, msg.StreamID, msg.Channel.ID, msg.Text, msg.Time).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
