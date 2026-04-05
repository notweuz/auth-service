package model

import "time"

type User struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	Username        string    `gorm:"not null;unique"`
	Password        string    `gorm:"not null"`
	PasswordVersion uint64    `gorm:"not null;default:1"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
