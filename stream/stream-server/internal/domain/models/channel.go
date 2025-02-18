package models

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID            uuid.UUID `json:"id,omitempty"`
	UserID        uuid.UUID `json:"user_id"`
	Live          bool      `json:"live"`
	ActiveStreamID bool      `json:"active_stream"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
