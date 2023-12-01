package services

import (
	"diploma/internal/api/repositories"
	"diploma/internal/models"
)

type AuthService struct {
	repository repositories.AuthRepository
}

func NewAuthService(r repositories.AuthRepository) AuthService {
	return AuthService{
		repository: r,
	}
}

func (s *AuthService) Get(login string) (*models.User, error) {
	return s.repository.Get(login)
}

func (s *AuthService) Register(user *models.User) error {
	return s.repository.Register(user)
}
