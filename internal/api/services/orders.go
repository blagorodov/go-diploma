package services

import (
	"diploma/internal/api/repositories"
	"diploma/internal/models"
)

type OrdersService struct {
	repository repositories.OrdersRepository
}

func NewOrdersService(r repositories.OrdersRepository) OrdersService {
	return OrdersService{
		repository: r,
	}
}

func (s *OrdersService) Add(number string) (*models.Order, error) {
	return s.repository.Add(number)
}

func (s *OrdersService) List() ([]*models.Order, error) {
	return s.repository.List()
}
