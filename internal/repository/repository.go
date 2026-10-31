package repository

import (
	"context"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User UserRepository
}

type UserRepository interface {
	Create(ctx context.Context, email string, hashPassword string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	q := sqlc.New(pool)

	return &Repositories{User: NewUserRepository(q)}
}
