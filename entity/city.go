package entity

import "gorm.io/gorm"

type City struct {
	gorm.Model
	ID    uint
	Name  string `gorm:"unique"`
	Store []Store
	User  []User
}
