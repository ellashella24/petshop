package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID      uint
	Name    string `gorm:"unique"`
	Product []Product
}
