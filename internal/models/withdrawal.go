package models

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Withdrawal struct {
	ID          int       `gorm:"primary_key;auto_increment" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	OrderNumber string    `gorm:"not null;size:500" json:"order_number"`
	Amount      int       `gorm:"not null;default:0" json:"amount"`
	ProcessedAt time.Time `gorm:"null" json:"processed_at"`
}

func (w *Withdrawal) TableName() string {
	return "withdrawals"
}

func (w *Withdrawal) ResponseMap() gin.H {
	resp := make(gin.H)
	resp["id"] = w.ID
	resp["user_id"] = w.UserID
	resp["order_number"] = w.OrderNumber
	resp["amount"] = w.Amount
	resp["processed_at"] = w.ProcessedAt
	return resp
}
