package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Title     string `json:"title"`
	Author    uint   `json:"author"`
}
