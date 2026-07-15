package repository

import (
	"cinerdle-back/internal/domain"
	"database/sql"

	"github.com/gofrs/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}

}
func (pst *PostgresUserRepository) Create(user *domain.User) error {
	row := pst.db.QueryRow("INSERT INTO users (email, password_hash, nickname, is_verified, wins, losses) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at ", user.Email, user.PasswordHash, user.Nickname, user.IsVerified, user.Wins, user.Losses)
	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return err
}

func (pst *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	row := pst.db.QueryRow("SELECT id, email, password_hash, nickname, is_verified, wins, losses, created_at, updated_at FROM users WHERE email = $1", email)
	user := &domain.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Nickname, &user.IsVerified, &user.Wins, &user.Losses, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (pst *PostgresUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	row := pst.db.QueryRow("SELECT id, email, password_hash, nickname, is_verified, wins, losses, created_at, updated_at FROM users WHERE id = $1", id)
	user := &domain.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Nickname, &user.IsVerified, &user.Wins, &user.Losses, &user.CreatedAt, &user.UpdatedAt,
	)
	return user, err
}

func (pst *PostgresUserRepository) Update(user *domain.User) error {
	_, err := pst.db.Exec("UPDATE users SET nickname=$1, is_verified=$2, wins=$3, losses=$4 WHERE id=$5", user.Nickname, user.IsVerified, user.Wins, user.Losses, user.ID)
	return err

}
