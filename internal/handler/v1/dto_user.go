package v1

import (
	"tgiintiz/internal/models"
	"time"
)

type userResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(u *models.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		RoleID:    int64(u.RoleID),
		CreatedAt: u.CreatedAt,
	}
}

type tokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newTokenPairResponse(tokens *models.TokenPair) tokenPairResponse {
	return tokenPairResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
