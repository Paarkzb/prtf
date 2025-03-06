package models

import (
	"time"

	"github.com/google/uuid"
)

type Recording struct {
	ID          uuid.UUID      `json:"id,omitempty"`
	ChannelId   uuid.UUID      `json:"channel_id,omitempty"`
	ChannelName string         `json:"channel_name,omitempty"`
	Path        *string        `json:"path,omitempty"`
	Date        time.Time      `json:"date,omitempty"`
	Poster      *string        `json:"poster,omitempty"`
	Duration    *time.Duration `json:"duration,omitempty"`
}
