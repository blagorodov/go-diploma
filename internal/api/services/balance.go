package services

import (
	"diploma/internal/api/repositories"
	"diploma/internal/auth"
	"diploma/internal/errs"
	"diploma/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
)

type BalanceService struct {
	balanceRepository repositories.BalanceRepository
	orderRepository   repositories.OrdersRepository
}

func NewBalanceService(balanceRepository repositories.BalanceRepository, orderRepository repositories.OrdersRepository) BalanceService {
	return BalanceService{
		balanceRepository: balanceRepository,
		orderRepository:   orderRepository,
	}
}

func (s *BalanceService) Get(ctx *gin.Context) (float64, float64, error) {
	userID, _ := auth.GetID(ctx)
	balance, err := s.balanceRepository.GetBalance(userID)
	if err != nil {
		return 0, 0, err
	}
	withdrawn, err := s.balanceRepository.GetWithdrawn(userID)
	if err != nil {
		return 0, 0, err
	}
	return float64(balance) / 100, float64(withdrawn) / 100, err
}

func (s *BalanceService) Withdraw(ctx *gin.Context, orderNumber string, sum float64) error {
	userID, _ := auth.GetID(ctx)
	// check order exists
	_, err := s.orderRepository.Get(orderNumber)
	if errors.Is(err, errs.ErrOrderNotFound) {
		return err
	}
	// check user balance
	balance, _, err := s.Get(ctx)
	if err != nil {
		return err
	}
	if balance < sum {
		return errs.ErrBalanceNotEnoughFunds
	}
	err = s.balanceRepository.Withdraw(orderNumber, userID, int(sum*100))
	if err != nil {
		return err
	}
	return err
}

func (s *BalanceService) Withdrawals(ctx *gin.Context) ([]*models.Withdrawal, error) {
	userID, _ := auth.GetID(ctx)
	return s.balanceRepository.Withdrawals(userID)
}
