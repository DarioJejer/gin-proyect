package models

import "gorm.io/gorm"

type House struct {
	gorm.Model
	Address string `json:"address"`
	Owner   uint   `json:"owner"`
}
