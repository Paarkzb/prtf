package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id", omitempty`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	PassHash []byte    `json:"pass_hash", omitempty`
}
