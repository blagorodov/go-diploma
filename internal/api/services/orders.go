package services

import (
	"context"
	"diploma/internal/api/clients"
	"diploma/internal/api/repositories"
	"diploma/internal/auth"
	"diploma/internal/errs"
	"diploma/internal/logger"
	"diploma/internal/models"
	"diploma/internal/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type OrdersService struct {
	ordersRepository repositories.OrdersRepository
	accrual          clients.AccrualClient
}

func NewOrdersService(or repositories.OrdersRepository, a clients.AccrualClient) OrdersService {
	return OrdersService{
		ordersRepository: or,
		accrual:          a,
	}
}

func (s *OrdersService) Add(ctx *gin.Context, number string) error {
	if !utils.ValidateLuhn(number) {
		return errs.ErrOrderNumberFormat
	}

	order, err := s.ordersRepository.Get(number)
	userID, _ := auth.GetID(ctx)

	if errors.Is(err, errs.ErrOrderNotFound) {
		return s.ordersRepository.Add(number, userID)
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
	return s.ordersRepository.ListAll(userID)
}

func (s *OrdersService) Poll(ctx context.Context, order *models.Order) (bool, error) {
	response, err := s.accrual.GetOrderInfo(ctx, order.Number)
	if err != nil {
		return false, fmt.Errorf("get order info: %w", err)
	}
	if response == nil {
		return false, nil
	}

	if order.Status == response.Status {
		return false, nil
	}

	order.Status = response.Status
	order.Accrual = (int)(response.Accrual * 100)

	if err = s.ordersRepository.Charge(order); err != nil {
		return false, err
	}

	return true, nil
}

func (s *OrdersService) RunPollingStatuses(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.pollStatuses(ctx); err != nil && !errors.Is(err, errs.ErrNoAccrual) {
				return err
			}
		}
	}
}

func (s *OrdersService) pollStatuses(ctx context.Context) error {
	if err := s.accrual.CanMakeRequest(); err != nil {
		logger.Log("Client cannot make a request")
		return nil
	}

	orders, err := s.ordersRepository.ListPending()
	if err != nil {
		return err
	}

	for _, order := range orders {
		updated, err := s.Poll(ctx, order)
		if err != nil {
			return err
		}
		if updated {
			logger.Log(fmt.Sprintf("Order %d updated", order.ID))
		}
	}

	return nil
}
