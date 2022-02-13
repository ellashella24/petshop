package cart

type CartRequest struct {
	UserID    uint `json:"userID"`
	ProductID uint `json:"productID"`
	Quantity  int  `json:"quantity"`
}
type CartTransactionRequest struct {
	ProductID int `json:"productID"`
}
