package service

import (
	"context"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
	"tgiintiz/internal/service/auth"
	"tgiintiz/internal/service/survey"
	"tgiintiz/internal/service/user"
)

type Services struct {
	Auth   AuthService
	User   UserService
	Survey SurveyService
}

type AuthService interface {
	Register(ctx context.Context, in models.RegisterInput) (*models.User, error)
	Login(ctx context.Context, in models.LoginInput) (*models.TokenPair, error)
	Authenticate(ctx context.Context, token string) (*models.TokenClaims, error)
	Refresh(ctx context.Context, refreshToken string) (*models.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type SurveyService interface {
	Create(ctx context.Context, userID int64, in models.SurveyRequestInput) (*models.SurveyRequest, error)
	GetByID(ctx context.Context, id int64, userID int64) (*models.SurveyRequest, error)
	ListByUserID(ctx context.Context, userID int64) ([]*models.SurveyRequest, error)
	ListAll(ctx context.Context) ([]*models.SurveyRequest, error)
	ListByStatus(ctx context.Context, status string) ([]*models.SurveyRequest, error)
	UpdateStatus(ctx context.Context, id int64, status string) (*models.SurveyRequest, error)
	Delete(ctx context.Context, id int64) error
}

func NewServices(repos *repository.Repositories, tokenCfg auth.TokenConfig) *Services {
	authService := auth.New(repos.User, repos.Session, tokenCfg)
	userService := user.New(repos.User)
	surveyService := survey.New(repos.Survey)

	return &Services{
		Auth:   authService,
		User:   userService,
		Survey: surveyService,
	}
}
