package repository

import (
	"cinerdle-back/internal/domain"
	"database/sql"
)

type PostgresMovieRepository struct {
	db *sql.DB
}

func NewPostgresMovieRepository(db *sql.DB) *PostgresMovieRepository {
	return &PostgresMovieRepository{
		db: db,
	}
}

func (pstm *PostgresMovieRepository) FindByID(id int) (*domain.Movie, error) {
	row := pstm.db.QueryRow("SELECT id,title, orgiginal_title, poster_path, release_date, overview, popularity, vote_average, created_at, updated_at FROM movies WHERE id=$1", id)
	movie := &domain.Movie{}
	err := row.Scan(
		&movie.ID, &movie.Title, &movie.OriginalTitle, &movie.PosterPath, &movie.ReleaseDate, &movie.Overview, &movie.Popularity, &movie.VoteAverage, &movie.CreatedAt, &movie.UpdatedAt,
	)
	return movie, err
}

func (pstm *PostgresMovieRepository) Save(movie *domain.Movie) error {
	row := pstm.db.QueryRow("INSERT INTO movies(id, title, orgiginal_title, poster_path, release_date, overview, popularity, vote_average) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at", movie.ID, movie.Title, movie.OriginalTitle, movie.PosterPath, movie.ReleaseDate, movie.Overview, movie.Popularity, movie.VoteAverage)
	err := row.Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	return err
}
