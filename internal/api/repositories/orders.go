package repositories

import (
	"diploma/internal/drivers"
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
func (r *OrdersRepository) Add(number string) (*models.Order, error) {
	return nil, nil
}

func (r *OrdersRepository) List() ([]*models.Order, error) {
	return nil, nil
}
