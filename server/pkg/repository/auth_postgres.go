package repository

import (
	"fmt"
	"prtf"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user prtf.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES($1, $2, $3) RETURNING id", userTable)

	err := r.db.QueryRow(context.Background(), query, user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (prtf.User, error) {
	var user prtf.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)

	err := r.db.QueryRow(context.Background(), query, username, password).Scan(&user.Id)

	return user, err
}

func (r *AuthPostgres) GetUserById(id uuid.UUID) (prtf.UserResponse, error) {
	var user prtf.UserResponse

	query := fmt.Sprintf("SELECT id, name, username FROM %s WHERE id=$1", userTable)

	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Name, &user.Username)

	return user, err
}
