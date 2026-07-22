package usecases

import (
	"cinerdle-back/internal/adapters/tmdb"
	"cinerdle-back/internal/domain"
	apperrors "cinerdle-back/pkg/errors"
	"time"
)

type MovieUseCase struct {
	tmdb            *tmdb.TMDbClient
	movieRepository domain.MovieRepository
}

func NewMovieUseCase(tmdb *tmdb.TMDbClient, movieRepository domain.MovieRepository) *MovieUseCase {
	return &MovieUseCase{
		tmdb:            tmdb,
		movieRepository: movieRepository,
	}
}

func (muc *MovieUseCase) SearchMovie(query string) ([]domain.Movie, error) {
	movie, err := muc.tmdb.SearchMovie(query)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	dmMovie := []domain.Movie{}
	for _, sortMovie := range movie {

		strTime := sortMovie.ReleaseDate
		parsedTime, err := time.Parse("2006-01-02", strTime)
		if err != nil {
			continue
		}
		dmMovie = append(dmMovie, domain.Movie{
			ID:          sortMovie.ID,
			Title:       sortMovie.Title,
			PosterPath:  sortMovie.PosterPath,
			ReleaseDate: parsedTime,
			Popularity:  sortMovie.Popularity,
			Overview:    sortMovie.Overview,
			VoteAverage: sortMovie.VoteAverage,
		})
	}
	return dmMovie, nil
}

func (muc *MovieUseCase) GetMovieCredits(movieID int) ([]domain.Person, error) {
	person, err := muc.tmdb.GetMovieCredits(movieID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	dmPerson := []domain.Person{}
	for _, sortCast := range person.Cast {

		dmPerson = append(dmPerson, domain.Person{
			ID:          sortCast.ID,
			Name:        sortCast.Name,
			ProfilePath: sortCast.ProfilePath,
		})

	}
	return dmPerson, err

}

func (muc *MovieUseCase) ValidateConnection(movieID1, movieID2 int) (bool, error) {
	firstMovie, err := muc.tmdb.GetMovieCredits(movieID1)

	if err != nil {
		return bool(false), apperrors.ErrInternalServer
	}

	secondMovie, err := muc.tmdb.GetMovieCredits(movieID2)
	if err != nil {
		return bool(false), apperrors.ErrInternalServer
	}

	for _, firstMovies := range firstMovie.Cast {
		for _, secondMovies := range secondMovie.Cast {
			if firstMovies == secondMovies {
				return bool(true), nil
			}
		}
	}
	return bool(false), err
}
