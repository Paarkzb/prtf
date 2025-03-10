package models

import (
	"github.com/google/uuid"
)

type Channel struct {
	ID          uuid.UUID `json:"id,omitempty"`
	ChannelName string    `json:"channel_name"`
	Icon        *string   `json:"icon,omitempty"`
}
