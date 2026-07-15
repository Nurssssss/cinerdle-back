package repository

import (
	"cinerdle-back/internal/domain"
	"database/sql"
)

type PostgresPersonRepository struct {
	db *sql.DB
}

func NewPostgresPersonRepository(db *sql.DB) *PostgresPersonRepository {
	return &PostgresPersonRepository{
		db: db,
	}
}

func (pstp *PostgresPersonRepository) FindByID(id int) (*domain.Person, error) {
	row := pstp.db.QueryRow("SELECT id, name, profile_path, birthday, deathday, biography, popularity, created_at, updated_at FROM persons WHERE id = $1", id)
	person := &domain.Person{}
	err := row.Scan(
		&person.ID,
		&person.Name,
		&person.ProfilePath,
		&person.Birthday,
		&person.DeathDay,
		&person.Biography,
		&person.Popularity,
		&person.CreatedAt,
		&person.UpdatedAt,
	)
	return person, err
}

func (pstp *PostgresPersonRepository) Save(person *domain.Person) error {
	row := pstp.db.QueryRow("INSERT INTO persons (id, name, profile_path, birthday, deathday, biography, popularity) VALUES($1, $2, $3,$4, $5, $6, $7) RETURNING id,created_at, updated_at", person.ID, person.Name, person.ProfilePath, person.Birthday, person.DeathDay, person.Biography, person.Popularity)

	err := row.Scan(
		&person.ID,
		&person.CreatedAt,
		&person.UpdatedAt,
	)
	return err
}
