package repository

import (
	"cinerdle-back/internal/domain"
	"database/sql"

	"github.com/gofrs/uuid"
)

type PostgresEmailVerificationRepository struct {
	db *sql.DB
}

func NewPostgreEmailVerificationRepository(db *sql.DB) *PostgresEmailVerificationRepository {
	return &PostgresEmailVerificationRepository{
		db: db,
	}
}

func (pstEm *PostgresEmailVerificationRepository) Create(email *domain.EmailVerification) error {
	_, err := pstEm.db.Exec("INSERT INTO email_verifications (user_id, code, expiration_at ) VALUES ($1, $2,$3 )", email.UserID, email.Code, email.ExpirationAt)
	return err

}

func (pstEm *PostgresEmailVerificationRepository) FindByCode(code string) (*domain.EmailVerification, error) {
	row := pstEm.db.QueryRow("SELECT id, user_id, code, expiration_at, created_at FROM email_verifications WHERE code = $1", code)
	email := &domain.EmailVerification{}
	err := row.Scan(
		&email.ID, &email.UserID, &email.Code, &email.ExpirationAt, &email.CreatedAt,
	)
	return email, err
}

func (pstEm *PostgresEmailVerificationRepository) DeleteByID(id uuid.UUID) error {
	_, err := pstEm.db.Exec("DELETE FROM email_verifications WHERE id= $1", id)
	return err
}
