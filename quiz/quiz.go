package prtf

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	Id          uuid.UUID  `json:"id" db:"id"`
	RfUserId    uuid.UUID  `json:"rf_user_id" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type SaveQuizInput struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type UpdateQuizInput struct {
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	Questions   *[]Question `json:"questions"`
}

type QuizResponse struct {
	Id          uuid.UUID    `json:"id"`
	User        UserResponse `json:"user"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Questions   []Question   `json:"questions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Question struct {
	Id     uuid.UUID `json:"id"`
	Index  int       `json:"index"`
	Title  string    `json:"title" binding:"required"`
	Answer string    `json:"answer" binding:"required"`
}

func (i UpdateQuizInput) Validate() error {
	if i.Name == nil && i.Description == nil && i.Questions != nil {
		return errors.New("update fields has no value")
	}

	return nil
}
