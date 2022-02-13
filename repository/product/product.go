package product

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
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

	if err != nil || len(products) == 0 {
		return products, errors.New("product not found")
	}

	return products, nil
}

func (gr *productRepository) GetProductByID(productID int) (entity.Product, error) {
	product := entity.Product{}

	err := gr.db.Where("id = ? AND category_id > 1", productID).Find(&product).Error

	if err != nil || product.ID == 0 {
		return product, errors.New("product not found")
	}

	return product, nil
}

func (gr *productRepository) GetProductByStoreID(storeID int) ([]entity.Product, error) {
	product := []entity.Product{}

	err := gr.db.Where("store_id = ?", storeID).Find(&product).Error

	if err != nil || len(product) == 0 {
		return product, errors.New("product not found")
	}

	return product, nil
}

func (gr *productRepository) GetStockHistory(productID int) ([]entity.StockHistory, error) {
	stock := []entity.StockHistory{}

	err := gr.db.Where("product_id = ?", productID).Find(&stock).Error

	if err != nil || len(stock) == 0 {
		return stock, errors.New("product history not found")
	}

	return stock, nil
}

func (gr *productRepository) CreateProduct(userID int, newProduct entity.Product) (entity.Product, error) {
	// var store entity.Store
	var stock entity.StockHistory
	// err := gr.db.Where("user_id = ? and id = ?", userID, newProduct.StoreID).First(&store).Error

	// if err != nil || store.ID == 0 {
	// 	return newProduct, errors.New("store not found")
	// }

	err := gr.db.Save(&newProduct).Error
	if err != nil {
		return newProduct, err
	}

	stock = entity.StockHistory{
		ProductID: newProduct.ID,
		Stock:     newProduct.Stock,
	}

	gr.db.Save(&stock)
	// err = gr.db.Save(&stock).Error
	// if err != nil {
	// 	return newProduct, err
	// }

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

		gr.db.Save(&stock)

		// err = gr.db.Save(&stock).Error
		// if err != nil {
		// 	return updatedProduct, err
		// }
	}

	gr.db.Model(&product).Updates(updatedProduct)

	return product, nil
}

func (gr *productRepository) DeleteProduct(productID int) (entity.Product, error) {
	product := entity.Product{}

	err := gr.db.Where("id = ?", productID).Find(&product).Error

	if err != nil || product.ID == 0 {
		return product, errors.New("product not found")
	}

	gr.db.Delete(&product)

	return product, nil
}
