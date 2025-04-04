package models

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID               uuid.UUID `json:"id,omitempty"`
	RfUserID         uuid.UUID `json:"rf_user_id"`
	ChannelName      string    `json:"channel_name"`
	Live             bool      `json:"live"`
	RfActiveStreamID uuid.UUID `json:"rf_active_stream_id,omitempty"`
	ChannelToken     string    `json:"channel_token,omitempty"`
	Icon             *string   `json:"icon,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
