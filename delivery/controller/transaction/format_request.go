package transaction

type TransactionRequest struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
	PetID     int `json:"petID"`
}

type TransactionDetailRequest struct {
	TransactionID uint
	ProductID     uint
	Quantity      uint
}

type InvoiceItemRequest struct {
	Name     string
	Price    int
	Quantity int
	Category string
}

type CallbackRequest struct {
	ExternalID    string `json:"external_id"`
	PaymentMethod string `json:"payment_method"`
	PaidAt        string `json:"paid_at"`
	Status        string `json:"status"`
}
