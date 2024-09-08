package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"prtf"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QuizPostgres struct {
	db *pgxpool.Pool
}

func NewQuizPostgres(db *pgxpool.Pool) *QuizPostgres {
	return &QuizPostgres{
		db: db,
	}
}

func (r *QuizPostgres) Save(userId uuid.UUID, quiz prtf.SaveQuizInput) (uuid.UUID, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	quesitons, err := json.Marshal(quiz.Questions)
	if err != nil {
		return uuid.Nil, err
	}

	var quizId uuid.UUID
	saveQuizQuery := fmt.Sprintf("INSERT INTO %s (rf_user_id, name, description, questions) VALUES($1, $2, $3, $4) RETURNING id", quizTable)
	err = tx.QueryRow(context.Background(), saveQuizQuery, userId, quiz.Name, quiz.Description, quesitons).Scan(&quizId)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return uuid.Nil, err
	}

	return quizId, tx.Commit(context.Background())
}

func (r *QuizPostgres) GetAll(userId uuid.UUID) ([]prtf.QuizResponse, error) {
	var quizes []prtf.QuizResponse

	query := fmt.Sprintf(`
		SELECT 
			quiz.id as quiz_id, 
			quiz.name as quiz_name, 
			quiz.description as quiz_description, 
			quiz.questions as quiz_question, 
			quiz.created_at as quiz_created_at,
			quiz.updated_at as quiz_updated_at,
			u.id as user_id,
			u.name as user_name,
			u.username as user_username
		FROM %s as quiz 
		LEFT JOIN public.user as u ON u.id = quiz.rf_user_id
		WHERE quiz.deleted=false
		ORDER BY quiz.created_at DESC`, quizTable)

	rows, err := r.db.Query(context.Background(), query)
	for rows.Next() {
		var quiz prtf.QuizResponse
		err = rows.Scan(&quiz.Id, &quiz.Name, &quiz.Description, &quiz.Questions, &quiz.CreatedAt, &quiz.UpdatedAt, &quiz.User.Id, &quiz.User.Name, &quiz.User.Username)
		if err != nil {
			return nil, err
		}

		quizes = append(quizes, quiz)
	}

	return quizes, err
}

func (r *QuizPostgres) GetById(userId, quizId uuid.UUID) (prtf.QuizResponse, error) {
	var quiz prtf.QuizResponse

	quizQuery := fmt.Sprintf(`
		SELECT 
			quiz.id, quiz.name, quiz.description, quiz.questions, quiz.created_at, quiz.updated_at, u.id, u.name, u.username
		FROM %s as quiz
		LEFT JOIN public.user as u on u.id=quiz.rf_user_id 
		WHERE quiz.id=$1`, quizTable)

	err := r.db.QueryRow(context.Background(), quizQuery, quizId).Scan(&quiz.Id, &quiz.Name, &quiz.Description, &quiz.Questions, &quiz.CreatedAt, &quiz.UpdatedAt, &quiz.User.Id, &quiz.User.Name, &quiz.User.Username)
	if err != nil {
		return quiz, err
	}

	return quiz, err
}

func (r *QuizPostgres) DeleteById(userId, quizId uuid.UUID) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}

	quizQuery := fmt.Sprintf("UPDATE %s as q SET deleted=true WHERE q.id=$1", quizTable)
	_, err = tx.Exec(context.Background(), quizQuery, quizId)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}

func (r *QuizPostgres) Update(userId, quizId uuid.UUID, input prtf.UpdateQuizInput) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}

	quizSetValues := make([]string, 0)
	quizArgs := make([]interface{}, 0)
	quizArgId := 1

	if input.Name != nil {
		quizSetValues = append(quizSetValues, fmt.Sprintf("name=$%d", quizArgId))
		quizArgs = append(quizArgs, *input.Name)
		quizArgId++
	}

	if input.Description != nil {
		quizSetValues = append(quizSetValues, fmt.Sprintf("description=$%d", quizArgId))
		quizArgs = append(quizArgs, *input.Description)
		quizArgId++
	}

	if input.Questions != nil {
		quizSetValues = append(quizSetValues, fmt.Sprintf("questions=$%d", quizArgId))
		quizArgs = append(quizArgs, *input.Questions)
		quizArgId++
	}

	quizSetValues = append(quizSetValues, fmt.Sprintf("updated_at=$%d", quizArgId))
	quizArgs = append(quizArgs, time.Now())
	quizArgId++

	setQuizQuery := strings.Join(quizSetValues, ", ")

	quizQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", quizTable, setQuizQuery, quizArgId)
	quizArgs = append(quizArgs, quizId)

	logrus.Debugf("Update query: %s", quizQuery)

	_, err = tx.Exec(context.Background(), quizQuery, quizArgs...)
	if err != nil {
		logrus.Info("JERE")
		_ = tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}
