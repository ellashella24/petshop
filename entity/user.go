package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID          uint
	Name        string
	Email       string `gorm:"unique"`
	Password    string
	CityID      uint
	Pet         []Pet
	Store       Store
	Role        string `gorm:"default:user"`
	Transaction []Transaction
	Cart        []Cart
}
