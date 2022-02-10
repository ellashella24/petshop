package product

import (
	"errors"
	"gorm.io/gorm"
	"petshop/entity"
)

type Product interface {
	GetAllProduct() ([]entity.Product, error)
	GetProductByID(productID int) (entity.Product, error)
	CreateProduct(userID int, product entity.Product) (entity.Product, error)
	UpdateProduct(productID int, product entity.Product) (entity.Product, error)
	DeleteProduct(productID int) (entity.Product, error)
	GetProductByStoreID(storeID int) ([]entity.Product, error)
	GetStockHistory(productID int) ([]entity.StockHistory, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db}
}

func (gr *productRepository) GetAllProduct() ([]entity.Product, error) {
	products := []entity.Product{}

	err := gr.db.Where("category_id > 1").Find(&products).Error

	if err != nil {
		return products, nil
	}

	return products, nil
}

func (gr *productRepository) GetProductByID(productID int) (entity.Product, error) {
	product := entity.Product{}

	err := gr.db.Where("id = ? AND id > 1", productID).Find(&product).Error

	if err != nil {
		return product, nil
	}

	return product, nil
}

func (gr *productRepository) GetProductByStoreID(storeID int) ([]entity.Product, error) {
	product := []entity.Product{}

	err := gr.db.Where("store_id = ?", storeID).Find(&product).Error

	if err != nil || len(product) == 0 {
		return product, err
	}

	return product, nil
}

func (gr *productRepository) GetStockHistory(productID int) ([]entity.StockHistory, error) {
	stock := []entity.StockHistory{}

	err := gr.db.Where("product_id = ?", productID).Find(&stock).Error

	if err != nil || len(stock) == 0 {
		return stock, err
	}

	return stock, nil
}

func (gr *productRepository) CreateProduct(userID int, newProduct entity.Product) (entity.Product, error) {
	var store entity.Store
	var stock entity.StockHistory
	err := gr.db.Where("user_id = ? and id = ?", userID, newProduct.StoreID).First(&store).Error

	if err != nil {
		return newProduct, err
	}

	err = gr.db.Save(&newProduct).Error
	if err != nil {
		return newProduct, err
	}

	stock = entity.StockHistory{
		ProductID: newProduct.ID,
		Stock:     newProduct.Stock,
	}

	err = gr.db.Save(&stock).Error
	if err != nil {
		return newProduct, err
	}

	return newProduct, nil
}

func (gr *productRepository) UpdateProduct(productID int, updatedProduct entity.Product) (entity.Product, error) {
	product := entity.Product{}
	var stock entity.StockHistory

	err := gr.db.Where("id = ?", productID).First(&product).Error

	if err != nil {
		return product, err
	}

	if product.CategoryID == 1 && product.Stock > 1 {
		return product, errors.New("Error")
	}

	if updatedProduct.Stock != 0 {
		updatedProduct.Stock = product.Stock + updatedProduct.Stock
		stock = entity.StockHistory{
			ProductID: product.ID,
			Stock:     updatedProduct.Stock,
		}
		err = gr.db.Save(&stock).Error
		if err != nil {
			return updatedProduct, err
		}
	}

	gr.db.Model(&product).Updates(updatedProduct)

	return updatedProduct, nil
}

func (gr *productRepository) DeleteProduct(productID int) (entity.Product, error) {
	product := entity.Product{}

	err := gr.db.Where("id = ?", productID).Delete(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}
