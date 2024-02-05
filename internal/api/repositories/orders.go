package repositories

import (
	"diploma/internal/drivers"
	"diploma/internal/errs"
	"diploma/internal/models"
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

func (r *OrdersRepository) List(userID string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
