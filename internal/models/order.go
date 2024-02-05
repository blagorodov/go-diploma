package models

import (
	"time"
)

const NEW = "NEW"
const PROCESSING = "PROCESSING"
const PROCESSED = "PROCESSED"
const INVALID = "INVALID"

type Order struct {
	ID         int       `gorm:"primary_key;auto_increment"`
	UserID     string    `gorm:"not null;index"`
	Number     string    `gorm:"not null;size:500"`
	Accrual    int       `gorm:"index"`
	Status     string    `gorm:"not null"`
	UploadedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
