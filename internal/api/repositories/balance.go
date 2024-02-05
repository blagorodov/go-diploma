package repositories

import (
	"diploma/internal/drivers"
	"diploma/internal/errs"
	"diploma/internal/models"
	"gorm.io/gorm"
	"time"
)

type BalanceRepository struct {
	db drivers.Database
}

type NResult struct {
	N int
}

func NewBalanceRepository(db drivers.Database) BalanceRepository {
	return BalanceRepository{
		db: db,
	}
}

func (r *BalanceRepository) GetBalance(userID string) (int, error) {
	var user models.User
	result := r.db.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return 0, errs.ErrUserNotFound
	}
	return user.Balance, nil
}

func (r *BalanceRepository) GetWithdrawn(userID string) (int, error) {
	var n NResult
	result := r.db.DB.Table("withdrawals").Where("user_id = ?", userID).Select("sum(amount) as n").Scan(&n)
	if result.Error != nil {
		return 0, result.Error
	}
	return n.N, nil
}

func (r *BalanceRepository) Withdraw(orderNumber string, userID string, sum int) error {
	withdrawal := models.Withdrawal{
		UserID:      userID,
		OrderNumber: orderNumber,
		Amount:      sum,
		ProcessedAt: time.Now(),
	}
	return r.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&withdrawal).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance - ?", sum)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *BalanceRepository) Withdrawals(userID string) ([]*models.Withdrawal, error) {
	var withdrawals []*models.Withdrawal
	if err := r.db.DB.Where("user_id = ?", userID).Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return withdrawals, nil
}
