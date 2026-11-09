package repository

import (
	"context"
	"tgiintiz/internal/models"
	attachmentRepo "tgiintiz/internal/repository/attachment"
	"tgiintiz/internal/repository/sqlc"
	surveyRepo "tgiintiz/internal/repository/survey"
	userRepo "tgiintiz/internal/repository/user"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User       UserRepository
	Survey     SurveyRepository
	Attachment AttachmentRepository
	Session    SessionRepository
}

type UserRepository interface {
	Create(ctx context.Context, email string, hashPassword string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error
	Exists(ctx context.Context, sessionID string) (bool, error)
	Delete(ctx context.Context, sessionID string) error
}

type SurveyRepository interface {
	Create(ctx context.Context, userID int64, input models.SurveyRequestInput) (*models.SurveyRequest, error)
	GetByID(ctx context.Context, id int64) (*models.SurveyRequest, error)
	ListByUserID(ctx context.Context, userID int64, limit, offset int32) ([]*models.SurveyRequest, error)
	ListAll(ctx context.Context, limit, ofset int32) ([]*models.SurveyRequest, error)
	ListByStatus(ctx context.Context, status string, limit, ofset int32) ([]*models.SurveyRequest, error)
	UpdateStatus(ctx context.Context, id int64, status string) (*models.SurveyRequest, error)
	Delete(ctx context.Context, id int64) error
}

type AttachmentRepository interface {
	Create(ctx context.Context, att *models.RequestAttachment) (*models.RequestAttachment, error)
	ListByRequestID(ctx context.Context, requestID int64) ([]*models.RequestAttachment, error)
	GetByID(ctx context.Context, id int64) (*models.RequestAttachment, error)
	Delete(ctx context.Context, id int64) error
}

func NewRepositories(pool *pgxpool.Pool, redisClient *redis.Client) *Repositories {
	q := sqlc.New(pool)

	return &Repositories{
		User:       userRepo.New(q),
		Survey:     surveyRepo.New(q),
		Attachment: attachmentRepo.New(q),
		Session:    NewSessionRepository(redisClient),
	}
}
