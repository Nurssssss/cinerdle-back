package jwt

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JwtManagerImpl struct {
	secret string
}

func NewJwtManagerImpl(secret string) *JwtManagerImpl {
	return &JwtManagerImpl{
		secret: secret,
	}
}

func (njm *JwtManagerImpl) Generate(userID uuid.UUID) (string, error) {
	claims := jwtlib.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix()}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(njm.secret))
}

func (njm *JwtManagerImpl) Verify(token string) (uuid.UUID, error) {
	parsedToken, err := jwtlib.Parse(token, func(parsedToken *jwtlib.Token) (interface{}, error) {
		return []byte(njm.secret), nil
	})
	if err != nil {
		return uuid.UUID{}, err

	}
	claims, ok := parsedToken.Claims.(jwtlib.MapClaims)
	if !ok {
		err := errors.New("invalid token")

		return uuid.UUID{}, err
	}

	sub := claims["sub"]
	subStr, ok := sub.(string)

	if !ok {
		err := errors.New("invalid sub")
		return uuid.UUID{}, err

	}
	return uuid.FromString(subStr)
}
