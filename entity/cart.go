package entity

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID        uint
	UserID    uint
	ProductID uint
	Jumlah    uint
}
