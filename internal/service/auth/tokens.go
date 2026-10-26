package auth

import (
	"fmt"
	"strconv"
	"tgiintiz/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Type string `json:"type"`
	Role string `json:"role,omitempty"`
}

func (s *Service) newTokenPair(u *models.User) (*models.TokenPair, error) {
	now := time.Now()

	accessToken, err := s.sign(tokenClaims{
		RegisteredClaims: registeredClaims(u.ID, now, s.accessTTL),
		Type:             accessTokenType,
		Role:             u.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := s.sign(tokenClaims{
		RegisteredClaims: registeredClaims(u.ID, now, s.refreshTTL),
		Type:             refreshTokenType,
		Role:             "",
	})
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken}, nil
}

func registeredClaims(userID int64, now time.Time, ttl time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
}

func (s *Service) sign(claims tokenClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
