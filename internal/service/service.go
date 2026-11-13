package service

import (
	"context"
	"tgiintiz/internal/models"
	"tgiintiz/internal/repository"
	"tgiintiz/internal/service/auth"
	"tgiintiz/internal/service/catalog"
	"tgiintiz/internal/service/survey"
	"tgiintiz/internal/service/user"
	"tgiintiz/pkg/storage"
)

type Services struct {
	Auth    AuthService
	User    UserService
	Survey  SurveyService
	Catalog CatalogService
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
	GetByID(ctx context.Context, requestID int64, userID int64) (*models.SurveyRequest, error)
	ListByUserID(ctx context.Context, userID int64, p models.Pagination) ([]*models.SurveyRequest, error)
	ListAll(ctx context.Context, p models.Pagination) ([]*models.SurveyRequest, error)
	ListByStatus(ctx context.Context, status string, p models.Pagination) ([]*models.SurveyRequest, error)
	UpdateStatus(ctx context.Context, requestID int64, status string) (*models.SurveyRequest, error)
	Delete(ctx context.Context, requestID int64) error

	UploadAttachment(ctx context.Context, in models.UploadAttachmentInput) (*models.RequestAttachment, error)
	ListAttachments(ctx context.Context, requestID int64, claims models.TokenClaims) ([]*models.RequestAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID int64, claims models.TokenClaims) error

	AddService(ctx context.Context, requestID, serviceID int64, claims models.TokenClaims) error
	ListServices(ctx context.Context, requestID int64, claims models.TokenClaims) ([]*models.ServiceCard, error)
	RemoveServices(ctx context.Context, requestID int64, claims models.TokenClaims) error
}

type CatalogService interface {
	Create(ctx context.Context, in models.CreateServiceInput) (*models.Service, error)
	GetByID(ctx context.Context, serviceID int64) (*models.Service, error)
	GetBySlugAndLang(ctx context.Context, slug, lang string) (*models.Service, error)
	ListByLang(ctx context.Context, lang string) ([]*models.ServiceCard, error)
	Update(ctx context.Context, serviceID int64, in models.UpdateServiceInput) (*models.Service, error)
	Delete(ctx context.Context, serviceID int64) error
}

func NewServices(repos *repository.Repositories, tokenCfg auth.TokenConfig, fileStorage storage.FileStorage) *Services {
	authService := auth.New(repos.User, repos.Session, tokenCfg)
	userService := user.New(repos.User)
	surveyService := survey.New(repos.Survey, repos.Attachment, repos.Catalog, fileStorage)
	catalogService := catalog.New(repos.Catalog)

	return &Services{
		Auth:    authService,
		User:    userService,
		Survey:  surveyService,
		Catalog: catalogService,
	}
}
