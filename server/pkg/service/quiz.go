package service

import (
	"prtf"
	"prtf/pkg/repository"

	"github.com/google/uuid"
)

type QuizService struct {
	reps repository.Quiz
}

func newQuizService(repo repository.Quiz) *QuizService {
	return &QuizService{
		reps: repo,
	}
}

func (s *QuizService) Save(userId uuid.UUID, quiz prtf.SaveQuizInput) (uuid.UUID, error) {
	return s.reps.Save(userId, quiz)
}

func (s *QuizService) GetAll(userId uuid.UUID) ([]prtf.QuizResponse, error) {
	return s.reps.GetAll(userId)
}

func (s *QuizService) GetById(userId, quizId uuid.UUID) (prtf.QuizResponse, error) {
	return s.reps.GetById(userId, quizId)
}

func (s *QuizService) DeleteById(userId, quizId uuid.UUID) error {
	return s.reps.DeleteById(userId, quizId)
}

func (s *QuizService) Update(userId, quizId uuid.UUID, input prtf.UpdateQuizInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.reps.Update(userId, quizId, input)
}
