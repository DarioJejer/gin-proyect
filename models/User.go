package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	CompanyID uint    `json:"company_id"`
	Company   Company `json:"company"`
	Books     []Book  `json:"books" gorm:"foreignKey:Author"`
	House     *House  `json:"house" gorm:"foreignKey:Owner"` // Pointer to allow null value
}
