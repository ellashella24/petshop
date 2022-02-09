package entity

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	ID      uint
	Name    string
	UserID  uint `gorm:"uniqueIndex"`
	CityID  uint
	Product []Product
}
