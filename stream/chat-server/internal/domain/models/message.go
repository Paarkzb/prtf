package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	StreamID uuid.UUID `json:"stream_id"`
	Text     string    `json:"text"`
	Time     time.Time `json:"time"`
	Channel  Channel   `json:"channel"`
}
