package repository

import (
	"context"
	"errors"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueViolation = "23505"

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *userRepository {
	return &userRepository{q: q}
}

func toModel(u sqlc.User) *models.User {
	return &models.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		CreatedAt:    u.CreatedAt.Time,
	}
}

func (r *userRepository) Create(ctx context.Context, email string, hashPassword string) (*models.User, error) {
	args := sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: hashPassword,
		Role:         "client",
	}

	user, err := r.q.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return nil, models.ErrEmailTaken
		}
		return nil, err
	}

	return toModel(user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return toModel(user), nil
}
