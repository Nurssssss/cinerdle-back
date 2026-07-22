package repository

import (
	"cinerdle-back/internal/domain"
	"database/sql"

	"github.com/gofrs/uuid"
)

type PostgresGameSessionRepository struct {
	db *sql.DB
}

func NewPostgressGameSessionRepository(db *sql.DB) *PostgresGameSessionRepository {
	return &PostgresGameSessionRepository{
		db: db,
	}
}

func (psgsm *PostgresGameSessionRepository) CreateRoom(gameSession *domain.GameSession) error {
	row := psgsm.db.QueryRow("INSERT INTO game_sessions(id, player1_id, player2_id, current_turn_user_id, status, last_movie_id, player1_skips_left, player2_skips_left, player1_hints_left, player2_hints_left, turn_ends_at, started_at, ended_at, winner_id) VALUES($1,$2, $3, $4,$5, $6, $7, $8, $9, $10, $11, $12, $13,$14) RETURNING id, created_at, updated_at", gameSession.ID, gameSession.FirstPlayerID, gameSession.SecondPlayerID, gameSession.CurrentTurnUserID, gameSession.Status, gameSession.LastMovieID, gameSession.FirstPlayerskipsLeft, gameSession.SecondPlayerSkipsLeft, gameSession.FirsPlayerHintsLeft, gameSession.SecondPlayerHintsLeft, gameSession.TurnEndsAt, gameSession.StartedAt, gameSession.EndedAt, gameSession.WinnerId)
	err := row.Scan(
		&gameSession.ID,
		&gameSession.CreatedAt,
		&gameSession.UpdatedAt,
	)
	return err
}

func (psgsm *PostgresGameSessionRepository) FindGame(id uuid.UUID) (*domain.GameSession, error) {
	row := psgsm.db.QueryRow("SELECT id, player1_id, player2_id, current_turn_user_id, status, last_movie_id, player1_skips_left, player2_skips_left, player1_hints_left, player2_hints_left, turn_ends_at, started_at, ended_at, winner_id FROM game_sessions WHERE id = $1", id)
	gameSession := &domain.GameSession{}
	err := row.Scan(
		&gameSession.ID,
		&gameSession.FirstPlayerID,
		&gameSession.SecondPlayerID,
		&gameSession.CurrentTurnUserID,
		&gameSession.Status,
		&gameSession.LastMovieID,
		&gameSession.FirstPlayerskipsLeft,
		&gameSession.SecondPlayerSkipsLeft,
		&gameSession.FirsPlayerHintsLeft,
		&gameSession.SecondPlayerHintsLeft,
		&gameSession.TurnEndsAt,
		&gameSession.StartedAt,
		&gameSession.EndedAt,
		&gameSession.WinnerId,
	)
	return gameSession, err
}
