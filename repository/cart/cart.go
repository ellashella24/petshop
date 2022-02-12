package cart

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"petshop/entity"
)

type Cart interface {
	GetAll(userId int) ([]entity.Cart, error)
	Create(entity.Cart) (entity.Cart, error)
	Update(entity.Cart) (entity.Cart, error)
	Delete(userId int, productId int) (entity.Cart, error)
	CheckCart(entity.Cart) (entity.Cart, error)
	GetProductByID(productID int) (entity.Product, error)
	GetCategoryByID(categoryID int) (entity.Category, error)
	GetUserByID(userID int) (entity.User, error)
	UpdateStock(productID, stock int) error
	Transaction(newTransactions entity.Transaction) (entity.Transaction, error)
	TransactionDetail(newDetailTransactions entity.TransactionDetail) error
}
type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db}
}
func (cr *CartRepository) GetAll(userId int) ([]entity.Cart, error) {
	var carts []entity.Cart

	if err := cr.db.Where(
		"user_id = ?", userId,
	).Find(&carts).Error; err != nil || len(carts) == 0 {
		return nil, errors.New("error")
	}

	return carts, nil
}
func (cr *CartRepository) CheckCart(c entity.Cart) (entity.Cart, error) {
	var cart entity.Cart

	if err := cr.db.Where("user_id = ? AND product_id = ?", c.UserID, c.ProductID).First(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}
func (cr *CartRepository) Create(cart entity.Cart) (entity.Cart, error) {
	var product entity.Product

	cr.db.Where("id", cart.ProductID).First(&product)
	if product.CategoryID == 1 {
		return cart, errors.New("error grooming")
	}
	if err := cr.db.Create(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}
func (cr *CartRepository) Update(cart entity.Cart) (entity.Cart, error) {
	var c entity.Cart

	if err := cr.db.Model(&c).Where(
		"user_id = ? AND product_id = ?", cart.UserID, cart.ProductID,
	).First(&c).Updates(cart).Error; err != nil {
		return c, err
	}

	return c, nil
}
func (cr *CartRepository) Delete(userId int, productId int) (entity.Cart, error) {
	var cart entity.Cart

	fmt.Println(userId, productId)
	err := cr.db.Where("user_id = ? AND product_id = ?", userId, productId).Delete(&cart).Error
	if err != nil {
		return cart, err
	}
	return cart, nil
}
func (tr *CartRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	err := tr.db.Save(&newTransactions).Error
	if err != nil {
		return newTransactions, err
	}
	return newTransactions, nil
}
func (tr *CartRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) error {
	err := tr.db.Save(&newDetailTransactions).Error
	if err != nil {
		return err
	}
	return nil
}

//Helper cart transaction
func (tr *CartRepository) GetProductByID(productID int) (entity.Product, error) {
	var product entity.Product
	err := tr.db.Where("id = ?", productID).First(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}
func (tr *CartRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	var category entity.Category
	err := tr.db.Where("id = ?", categoryID).First(&category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}
func (tr *CartRepository) GetUserByID(userID int) (entity.User, error) {
	var user entity.User
	err := tr.db.Where("id = ?", userID).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
func (tr *CartRepository) UpdateStock(productID, stock int) error {
	var product entity.Product

	err := tr.db.Where("id = ?", productID).Model(&product).Update("stock", stock).Error

	if err != nil {
		return err
	}

	return nil
}
