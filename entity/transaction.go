package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID                uint
	UserID            uint
	InvoiceID         string
	PaymentMethod     string
	PaymentURL        string
	PaidAt            time.Time
	TotalPrice        int
	PaymentStatus     string
	TransactionDetail []TransactionDetail
}
