package entity

import "gorm.io/gorm"

type TransactionDetail struct {
	gorm.Model
	ID             uint
	TransactionID  uint
	ProductID      uint
	Quantity       int
	GroomingStatus GroomingStatus
}
