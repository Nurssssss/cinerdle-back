package domain

import (
	"time"
)

type Movie struct {
	ID            int
	Title         string
	OriginalTitle string
	PosterPath    string
	ReleaseDate   time.Time
	Overview      string
	Popularity    float64
	VoteAverage   float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type MovieRepository interface {
	FindByID(id int) (*Movie, error)
	Save(*Movie) error
}
