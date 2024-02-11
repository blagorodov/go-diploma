package repositories

import (
	"diploma/internal/drivers"
	"diploma/internal/errs"
	"diploma/internal/models"
	"gorm.io/gorm"
)

type OrdersRepository struct {
	db drivers.Database
}

func NewOrdersRepository(db drivers.Database) OrdersRepository {
	return OrdersRepository{
		db: db,
	}
}

func (r *OrdersRepository) Get(number string) (*models.Order, error) {
	var order models.Order
	result := r.db.DB.Where("number = ?", number).First(&order)
	if result.Error != nil {
		return nil, errs.ErrOrderNotFound
	}
	return &order, nil
}

func (r *OrdersRepository) Add(number string, userID string) error {
	result := r.db.DB.Create(&models.Order{
		UserID: userID,
		Status: models.NEW,
		Number: number,
	})
	return result.Error
}

func (r *OrdersRepository) Update(order *models.Order) error {
	result := r.db.DB.Save(order)
	return result.Error
}

func (r *OrdersRepository) ListAll(userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrdersRepository) ListPending() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.DB.
		Where("status in (?, ?)", models.NEW, models.PROCESSING).
		Find(&orders).Error; err != nil {

		return nil, err
	}
	return orders, nil
}

func (r *OrdersRepository) Charge(order *models.Order) error {
	return r.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).
			Where("id = ?", order.UserID).
			Update("balance", gorm.Expr("balance - ?", order.Accrual)).Error; err != nil {

			return err
		}
		return nil
	})
}
