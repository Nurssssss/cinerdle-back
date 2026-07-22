package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type GameMove struct {
	UserID    uuid.UUID
	ID        uuid.UUID
	SessionID uuid.UUID
	MovieID   int
	MoveType  string
	IsValid   bool
	CreatedAt time.Time
}
type GameMoveRepository interface {
}
