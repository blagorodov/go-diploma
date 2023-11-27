package auth

import (
	authrepositories "diploma/internal/api/repositories"
	"diploma/internal/models"
)

type Service struct {
	repository authrepositories.Repository
}

func NewService(r authrepositories.Repository) Service {
	return Service{
		repository: r,
	}
}

func (s *Service) Get(login string) (*models.User, error) {
	return s.repository.Get(login)
}

func (s *Service) Register(user *models.User) error {
	return s.repository.Register(user)
}
