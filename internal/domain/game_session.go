package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type GameSession struct {
	ID                    uuid.UUID
	FirstPlayerID         uuid.UUID
	SecondPlayerID        uuid.UUID
	CurrentTurnUserID     uuid.UUID
	LastMovieID           int
	Status                string
	FirstPlayerskipsLeft  int
	SecondPlayerSkipsLeft int
	FirsPlayerHintsLeft   int
	SecondPlayerHintsLeft int
	TurnEndsAt            time.Time
	StartedAt             time.Time
	EndedAt               time.Time
	WinnerId              uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type GameSessionRepository interface {
	CreateRoom(*GameSession) error
	FindGame(id uuid.UUID) (*GameSession, error)
}
