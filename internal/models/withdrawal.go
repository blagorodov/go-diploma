package models

import (
	"time"
)

type Withdrawal struct {
	ID          int       `gorm:"primary_key;auto_increment" json:"id"`
	UserID      string    `gorm:"not null;index" json:"user_id"`
	OrderNumber string    `gorm:"not null;size:500" json:"order_number"`
	Amount      int       `gorm:"not null;default:0" json:"amount"`
	ProcessedAt time.Time `gorm:"null" json:"processed_at"`
}
