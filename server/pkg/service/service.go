package service

import (
	"prtf"
	"prtf/pkg/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(user prtf.User) (uuid.UUID, error)
	GetUser(userId uuid.UUID) (prtf.UserResponse, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}

type Quiz interface {
	Save(userId uuid.UUID, quiz prtf.SaveQuizInput) (uuid.UUID, error)
	GetAll(userId uuid.UUID) ([]prtf.QuizResponse, error)
	GetById(userId, quizId uuid.UUID) (prtf.QuizResponse, error)
	DeleteById(userId, quizId uuid.UUID) error
	Update(userId, quizId uuid.UUID, input prtf.UpdateQuizInput) error
}

type Service struct {
	Authorization
	Quiz
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos.Authorization), Quiz: newQuizService(repos.Quiz)}
}
