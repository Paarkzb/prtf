package models

import "github.com/google/uuid"

type App struct {
	ID uuid.UUID
	Name string
	Secret string
}