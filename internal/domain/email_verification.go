package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type EmailVerification struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Code         string
	ExpirationAt time.Time
	CreatedAt    time.Time
}

type EmailVerificationRepository interface {
	Create(*EmailVerification) error
	FindByCode(code string) (*EmailVerification, error)
	DeleteByID(id uuid.UUID) error
}
