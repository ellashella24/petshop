package product

type ProductFormatResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
	ImageUrl string `json:"imageUrl"`
}
type StockFormatResponse struct {
	ID        uint `json:"id"`
	ProductID uint `json:"productID"`
	Stock     int  `json:"stock"`
}
