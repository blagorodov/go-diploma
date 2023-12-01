package models

import (
	"github.com/gin-gonic/gin"
	"time"
)

const NEW = "NEW"
const PROCESSING = "PROCESSING"
const PROCESSED = "PROCESSED"
const INVALID = "INVALID"

type Order struct {
	ID         int       `gorm:"primary_key;auto_increment" json:"id"`
	UserID     int       `gorm:"not null;index" json:"user_id"`
	Number     string    `gorm:"not null;size:500" json:"number"`
	Accrual    int       `gorm:"index" json:"accrual;omitempty"`
	Status     string    `gorm:"not null" json:"status"`
	UploadedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (o *Order) TableName() string {
	return "orders"
}

func (o *Order) ResponseMap() gin.H {
	resp := make(gin.H)
	resp["id"] = o.ID
	resp["user_id"] = o.UserID
	resp["number"] = o.Number
	resp["accrual"] = o.Accrual
	resp["status"] = o.Status
	resp["uploaded_at"] = o.UploadedAt
	resp["updated_at"] = o.UpdatedAt
	return resp
}
