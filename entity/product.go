package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID                uint
	Name              string
	Price             int
	Stock             int
	ImageURL          string
	StoreID           uint
	CategoryID        uint
	Category          Category
	Cart              []Cart
	TransactionDetail []TransactionDetail
	StockHistory      []StockHistory
}
