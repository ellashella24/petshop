package product

type CreateProductRequestFormat struct {
	Name       string `json:"name" form:"name" `
	Price      int    `json:"price" form:"price"`
	Stock      int    `json:"stock" form:"stock"`
	CategoryID uint   `json:"categoryID" form:"categoryid"`
}

type UpdateProductRequestFormat struct {
	Name     string `json:"name" form:"name" `
	Price    int    `json:"price" form:"price"`
	Stock    int    `json:"stock" form:"stock"`
	ImageURL string `json:"imageURL" form:"imageurl"`
}
