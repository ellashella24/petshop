package transaction

import "time"

type GetStoreTransactionResponse struct {
	Transaction       StoreTransactionResponse
	TransactionDetail []TransactionDetailResponse
}

type GetUserTransactionResponse struct {
	Transaction       TransactionResponse
	TransactionDetail []TransactionDetailResponse
}

type TransactionResponse struct {
	ID            uint      `json:"ID"`
	UserID        uint      `json:"userID"`
	InvoiceID     string    `json:"invoiceID"`
	PaymentMethod string    `json:"paymentMethod"`
	PaymentURL    string    `json:"paymentURL"`
	PaidAt        time.Time `json:"paidAt"`
	TotalPrice    int       `json:"totalPrice"`
	PaymentStatus string    `json:"paymentStatus"`
}

type StoreTransactionResponse struct {
	ID            uint   `json:"ID"`
	InvoiceID     string `json:"invoiceID"`
	PaymentStatus string `json:"paymentStatus"`
}

type TransactionDetailResponse struct {
	TransactionID uint `json:"transactionID"`
	ProductID     uint `json:"productID"`
	Quantity      int  `json:"quantity"`
}
