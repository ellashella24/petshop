package entity

import "gorm.io/gorm"

type StockHistory struct {
	gorm.Model
	ID        uint
	ProductID uint
	Stock     int
}
