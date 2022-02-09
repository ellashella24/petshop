package entity

import "gorm.io/gorm"

type GroomingStatus struct {
	gorm.Model
	ID                  uint
	PetID               uint
	Status              string
	TransactionDetailID uint `gorm:"uniqueIndex"`
}
