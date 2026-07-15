package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Nickname     string
	IsVerified   bool
	Wins         int
	Losses       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	Create(*User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uuid.UUID) (*User, error)
	Update(*User) error
}
