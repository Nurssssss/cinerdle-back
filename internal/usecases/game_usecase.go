package usecases

import (
	"cinerdle-back/internal/domain"
	apperrors "cinerdle-back/pkg/errors"

	"github.com/gofrs/uuid"
)

type GameUseCase struct {
	gameSessionRepository domain.GameSessionRepository
}

func NewGameSessionUseCase(gameSessionRepository domain.GameSessionRepository) *GameUseCase {
	return &GameUseCase{
		gameSessionRepository: gameSessionRepository,
	}
}

func (guc *GameUseCase) CreateGame(firstPlayer, secondPlayer uuid.UUID) (*domain.GameSession, error) {
	gameSession := domain.GameSession{
		FirstPlayerID:  firstPlayer,
		SecondPlayerID: secondPlayer,
		Status:         "waiting",
	}
	err := guc.gameSessionRepository.CreateRoom(&gameSession)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return &gameSession, nil

}
