package product

type CreateProductFormatResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
	ImageUrl string `json:"imageUrl"`
}

type UpdateFormatResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
	ImageUrl string `json:"imageUrl"`
}

type ProductFormatResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
	ImageUrl string `json:"imageUrl"`
	Category string `json:"category"`
}
type StockFormatResponse struct {
	ID        uint `json:"id"`
	ProductID uint `json:"productID"`
	Stock     int  `json:"stock"`
}
