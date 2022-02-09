package transaction

type TransactionRequest struct {
	ProductID int
	Quantity  int
	PetID     int
}

type TransactionDetailRequest struct {
	TransactionID uint
	ProductID     uint
	Quantity      uint
}
