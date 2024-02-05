package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	Login     string    `gorm:"not null;size:200;uniqueIndex" json:"login" binding:"required"`
	Password  string    `gorm:"size:200" json:"password" binding:"required"`
	Balance   int       `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err == nil {
		u.Password = string(bytes)
	}
	return err
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
