package auth

import (
	"context"
	"crypto/rand"
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

func (s *Service) newTokenPair(ctx context.Context, u *models.User) (*models.TokenPair, error) {
	sessionID := rand.Text()
	if err := s.sessionRepo.Create(ctx, sessionID, u.ID, s.refreshTTL); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	now := time.Now()

	accessToken, err := s.sign(tokenClaims{
		RegisteredClaims: registeredClaims(u.ID, sessionID, now, s.accessTTL),
		Type:             accessTokenType,
		Role:             u.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := s.sign(tokenClaims{
		RegisteredClaims: registeredClaims(u.ID, sessionID, now, s.refreshTTL),
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

func registeredClaims(userID int64, sessionID string, now time.Time, ttl time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ID:        sessionID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
}

func (s *Service) sign(claims tokenClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *Service) parse(tokenStr string, wantType string) (*tokenClaims, error) {
	claims := &tokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenStr, claims,
		func(*jwt.Token) (any, error) { return s.jwtSecret, nil },

		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", models.ErrInvalidToken, err)
	}

	if !parsedToken.Valid {
		return nil, models.ErrInvalidToken
	}

	if claims.Type != wantType {
		return nil, fmt.Errorf("%w: expected %s token, got %q", models.ErrInvalidToken, wantType, claims.Type)
	}

	return claims, nil
}

func (c *tokenClaims) userID() (int64, error) {
	userID, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid user id in subject", models.ErrInvalidToken)
	}

	return userID, nil
}
