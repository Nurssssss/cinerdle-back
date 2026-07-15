package usecases

import (
	"cinerdle-back/internal/adapters/email"
	"cinerdle-back/internal/domain"
	apperrors "cinerdle-back/pkg/errors"
	"cinerdle-back/pkg/jwt"
	"cinerdle-back/pkg/password"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
)

type AuthUseCase struct {
	userRepo              domain.UserRepository
	emailVerificationRepo domain.EmailVerificationRepository
	passwordHasher        password.PasswordHasher
	jwtManager            jwt.JwtManager
	emailVerification     email.EmailSender
}

func NewAuthCase(userRepo domain.UserRepository, emailVerificationRepo domain.EmailVerificationRepository, passwordHasher password.PasswordHasher, jwtManager jwt.JwtManager, emailVerification email.EmailSender) *AuthUseCase {
	return &AuthUseCase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		passwordHasher:        passwordHasher,
		jwtManager:            jwtManager,
		emailVerification:     emailVerification,
	}
}
func (ac *AuthUseCase) Register(ctx context.Context, email string, password string, nickname string) (*domain.User, error) {
	_, err := ac.userRepo.FindByEmail(email)
	if err == nil {
		return nil, apperrors.ErrEmailAlreadyExists
	}
	var hash string
	hash, err = ac.passwordHasher.Hash(password)
	if err != nil {
		log.Println(err)
		return nil, apperrors.ErrInternalServer
	}
	user := &domain.User{
		Email:        email,
		PasswordHash: hash,
		Nickname:     nickname,
	}
	err = ac.userRepo.Create(user)
	if err != nil {
		log.Println(err)
		return nil, apperrors.ErrInternalServer
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	verification := &domain.EmailVerification{
		UserID:       user.ID,
		Code:         code,
		ExpirationAt: time.Now().Add(time.Hour * 24),
	}
	err = ac.emailVerificationRepo.Create(verification)
	if err != nil {
		log.Println(err)
		return nil, apperrors.ErrInternalServer
	}
	err = ac.emailVerification.SendVerificationCode(user.Email, code)
	if err != nil {
		log.Println(err)
		return nil, apperrors.ErrInternalServer
	}
	return user, nil

}

func (ac *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := ac.userRepo.FindByEmail(email)
	if err != nil {
		return "", apperrors.ErrUserNotFound
	}
	ok := ac.passwordHasher.Check(password, user.PasswordHash)
	if !ok {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := ac.jwtManager.Generate(user.ID)
	if err != nil {
		log.Println(err)
		return "", apperrors.ErrInternalServer
	}
	return token, nil
}

func (ac *AuthUseCase) VerifyEmail(ctx context.Context, code string) error {
	verification, err := ac.emailVerificationRepo.FindByCode(code)
	if err != nil {
		return apperrors.ErrCodeVerififactionExpired
	}
	if time.Now().After(verification.ExpirationAt) {
		return apperrors.ErrCodeVerififactionExpired
	}
	userVerified, err := ac.userRepo.FindByID(verification.UserID)
	if userVerified == nil {
		return apperrors.ErrUserNotFound
	}

	userVerified.IsVerified = true
	err = ac.userRepo.Update(userVerified)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	err = ac.emailVerificationRepo.DeleteByID(verification.ID)
	if err != nil {
		return apperrors.ErrInternalServer
	}

	return nil
}

func (ac *AuthUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := ac.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	return user, nil

}
