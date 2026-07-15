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

func (muc *MovieUseCase) GetMovieCredits(movieID int) (*domain.Person, error) {
	person, err := muc.tmdb.GetMovieCredits(movieID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	tmdbCast := []tmdb.CastMember{}
	for _, sortCast := range person {

	}
}
