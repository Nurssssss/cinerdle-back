package jwt

import (
	"github.com/gofrs/uuid"
)

type JwtManager interface {
	Generate(userID uuid.UUID) (string, error)
	Verify(token string) (uuid.UUID, error)
}
