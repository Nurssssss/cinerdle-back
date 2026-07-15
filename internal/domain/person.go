package domain

import "time"

type Person struct {
	ID          int
	Name        string
	ProfilePath string
	Birthday    string
	DeathDay    time.Time
	Biography   string
	Popularity  float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PersonRepository interface {
	FindByID(id int) (*Person, error)
	Save(*Person) error
}
