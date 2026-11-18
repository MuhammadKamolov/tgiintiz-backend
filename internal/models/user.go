package models

import (
	"time"
)

type RoleID int64

const (
	RoleClient RoleID = 1
	RoleAdmin  RoleID = 2
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	RoleID       RoleID
	CreatedAt    time.Time
}

type RegisterInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
type TokenClaims struct {
	UserID int64
	RoleID RoleID
}

func (c TokenClaims) IsAdmin() bool {
	return c.RoleID == RoleAdmin
}
