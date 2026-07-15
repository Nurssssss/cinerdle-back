package http

import (
	"cinerdle-back/internal/domain"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type AuthHandler struct {
	authUseCase AuthUseCase
}
type AuthUseCase interface {
	Register(ctx context.Context, email string, password string, nickname string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	VerifyEmail(ctx context.Context, code string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

func NewAuthHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type registerResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Nickname   string    `json:"nickname"`
	IsVerified bool      `json:"is_verified"`
	Wins       int       `json:"wins"`
	Loses      int       `json:"loses"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type profileResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Nickname   string    `json:"nickname"`
	IsVerified bool      `json:"is_verified"`
	Wins       int       `json:"wins"`
	Loses      int       `json:"loses"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := ah.authUseCase.Register(r.Context(), req.Email, req.Password, req.Nickname)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	regResponse := &registerResponse{
		ID:         user.ID,
		Email:      user.Email,
		Nickname:   user.Nickname,
		IsVerified: user.IsVerified,
		Wins:       user.Wins,
		Loses:      user.Losses,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
	json.NewEncoder(w).Encode(regResponse)

}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := ah.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

type verifyRequest struct {
	Code string `json:"code"`
}

func (ah *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ah.authUseCase.VerifyEmail(r.Context(), req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

type getProfileRequest struct {
	UserID uuid.UUID `json:"userID"`
}

func (ah *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)
	user, err := ah.authUseCase.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	profileResponse := &profileResponse{
		ID:         user.ID,
		Email:      user.Email,
		Nickname:   user.Nickname,
		IsVerified: user.IsVerified,
		Wins:       user.Wins,
		Loses:      user.Losses,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
	json.NewEncoder(w).Encode(profileResponse)
}
