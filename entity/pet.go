package entity

import "gorm.io/gorm"

type Pet struct {
	gorm.Model
	ID             uint
	Name           string
	UserID         uint
	GroomingStatus []GroomingStatus
}
