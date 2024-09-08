package repository

import (
	"prtf"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(user prtf.User) (uuid.UUID, error)
	GetUser(username, password string) (prtf.User, error)
	GetUserById(userId uuid.UUID) (prtf.UserResponse, error)
}

type Quiz interface {
	Save(userId uuid.UUID, quiz prtf.SaveQuizInput) (uuid.UUID, error)
	GetAll(userId uuid.UUID) ([]prtf.QuizResponse, error)
	GetById(userId, quizId uuid.UUID) (prtf.QuizResponse, error)
	DeleteById(userId, quizId uuid.UUID) error
	Update(userId, quizId uuid.UUID, input prtf.UpdateQuizInput) error
}

type Repository struct {
	Authorization
	Quiz
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{Authorization: NewAuthPostgres(pool), Quiz: NewQuizPostgres(pool)}
}
