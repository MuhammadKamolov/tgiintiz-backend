package repository

import (
	"context"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User    UserRepository
	Session SessionRepository
}

type UserRepository interface {
	Create(ctx context.Context, email string, hashPassword string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

func NewRepositories(pool *pgxpool.Pool, redisClient *redis.Client) *Repositories {
	q := sqlc.New(pool)

	return &Repositories{
		User:    NewUserRepository(q),
		Session: NewSessionRepository(redisClient),
	}
}

type SessionRepository interface {
	Create(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error
	Exists(ctx context.Context, sessionID string) (bool, error)
	Delete(ctx context.Context, sessionID string) error
}
