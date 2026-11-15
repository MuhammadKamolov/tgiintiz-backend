package v1

import "tgiintiz/internal/models"

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (req registerRequest) toInput() models.RegisterInput {
	return models.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (req loginRequest) toInput() models.LoginInput {
	return models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
