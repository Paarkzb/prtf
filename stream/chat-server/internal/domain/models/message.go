package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	StreamChannelId uuid.UUID `json:"stream_channel_id"`
	StreamID  uuid.UUID `json:"stream_id"`
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
	Channel   Channel   `json:"channel"`
}
