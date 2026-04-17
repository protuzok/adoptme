package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrUnexpectedSigningMethod is returned when the JWT signing method is not expected.
var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

type Manager struct {
	secret   string
	duration time.Duration
}

func New(secret string, duration time.Duration) *Manager {
	return &Manager{
		secret:   secret,
		duration: duration,
	}
}

func (m *Manager) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.duration)),
	})

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("jwt - GenerateToken - token.SignedString: %w", err)
	}

	return tokenString, nil
}

func (m *Manager) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}

		return []byte(m.secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("jwt - ParseToken - jwtlib.Parse: %w", err)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("jwt - ParseToken - GetSubject: %w", err)
	}

	return sub, nil
}
