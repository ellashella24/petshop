package store

type StoreFormatResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	CityID uint   `json:"city_id"`
	UserID uint   `json:"user_id"`
}

type ListTransactionsStoreResponse struct {
	InvoiceID     string `json:"invoice_id"`
	UserID        uint   `json:"user_id"`
	TotalPrice    int    `json:"total_price"`
	PaymentStatus string `json:"payment_status"`
	PaidAt        string `json:"paid_at"`
	PaymentMethod string `json:"payment_method"`
}

type GroomingStatusResponse struct {
	ID     uint   `json:"id"`
	PetID  uint   `json:"pet_id"`
	Status string `json:"status"`
}
