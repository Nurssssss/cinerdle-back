package http

import (
	"cinerdle-back/internal/domain"
	"encoding/json"
	"net/http"
)

type MovieHandler struct {
	movieUseCase MovieUseCase
}
type MovieUseCase interface {
	SearchMovie(query string) ([]domain.Movie, error)
}

func NewMovieHandler(movieUseCase MovieUseCase) *MovieHandler {
	return &MovieHandler{
		movieUseCase: movieUseCase,
	}
}

func (mh *MovieHandler) SearchMovie(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	movie, err := mh.movieUseCase.SearchMovie(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movie)

}
