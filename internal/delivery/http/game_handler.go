package http

import (
	"cinerdle-back/internal/domain"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
)

type GameHandler struct {
	gameUseCase GameUseCase
}
type GameUseCase interface {
	CreateGame(firstPlayer, secondPlayer uuid.UUID) (*domain.GameSession, error)
}

func NewGameHandler(gameUseCase GameUseCase) *GameHandler {
	return &GameHandler{
		gameUseCase: gameUseCase,
	}
}

type gameRequest struct {
	FirstPlayer  uuid.UUID `json:"first_player"`
	SecondPlayer uuid.UUID `json:"second_player"`
}

func (guh *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game gameRequest
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameUseCase, err := guh.gameUseCase.CreateGame(game.FirstPlayer, game.SecondPlayer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameUseCase)

}
