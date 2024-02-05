package services

import (
	"diploma/internal/api/repositories"
	"diploma/internal/auth"
	"diploma/internal/errs"
	"diploma/internal/models"
	"diploma/internal/requests"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AuthService struct {
	repository repositories.AuthRepository
}

func NewAuthService(r repositories.AuthRepository) AuthService {
	return AuthService{
		repository: r,
	}
}

func (s *AuthService) Login(ctx *gin.Context, userRequest requests.User) error {
	userModel, err := s.repository.Get(userRequest.Login)
	if err != nil {
		return err
	}
	if !userModel.CheckPasswordHash(userRequest.Password) {
		return errs.ErrLoginOrPasswordNotFound
	}
	auth.SetID(ctx, strconv.Itoa(userModel.ID))
	return nil
}

func (s *AuthService) Register(ctx *gin.Context, userRequest requests.User) error {
	userModel := models.User{
		Login:    userRequest.Login,
		Password: userRequest.Password,
	}
	if err := userModel.HashPassword(); err != nil {
		return err
	}
	if err := s.repository.Register(&userModel); err != nil {
		return err
	}
	auth.SetID(ctx, strconv.Itoa(userModel.ID))
	return nil
}
