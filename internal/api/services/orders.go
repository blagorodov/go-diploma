package services

import (
	"diploma/internal/api/repositories"
	"diploma/internal/auth"
	"diploma/internal/errs"
	"diploma/internal/models"
	"diploma/internal/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

type OrdersService struct {
	repository repositories.OrdersRepository
}

func NewOrdersService(r repositories.OrdersRepository) OrdersService {
	return OrdersService{
		repository: r,
	}
}

func (s *OrdersService) Add(ctx *gin.Context, number string) error {
	if !utils.ValidateLuhn(number) {
		return errs.ErrOrderNumberFormat
	}
	order, err := s.repository.Get(number)
	userID, _ := auth.GetID(ctx)
	if errors.Is(err, errs.ErrOrderNotFound) {
		return s.repository.Add(number, userID)
	}
	if err == nil {
		if order.UserID != userID {
			return errs.ErrOrderOtherUserDuplicate
		}
		return errs.ErrOrderDuplicate
	}
	return nil
}

func (s *OrdersService) List(ctx *gin.Context) ([]*models.Order, error) {
	userID, _ := auth.GetID(ctx)
	return s.repository.List(userID)
}
