package transaction

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"petshop/entity"
)

type Transaction interface {
	Transaction(newTransactions entity.Transaction) (entity.Transaction, error)
	TransactionDetail(newDetailTransactions entity.TransactionDetail) (entity.TransactionDetail, error)
	GroomingStatusHelper(petID int, transactionDetailID uint) error
	PetValidator(petID int, userID int) error
	GetProductByID(productID int) (entity.Product, error)
	GetCategoryByID(categoryID int) (entity.Category, error)
	GetUserByID(userID int) (entity.User, error)
	Callback(callback entity.Transaction) error
	UpdateStock(productID, stock int) error
	GetAllUserTransaction(userID int) ([]entity.Transaction, error)
	GetAllStoreTransaction(userID int) ([]entity.TransactionDetail, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (tr *transactionRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	err := tr.db.Save(&newTransactions).Error
	if err != nil {
		return newTransactions, err
	}
	return newTransactions, nil
}

func (tr *transactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	err := tr.db.Save(&newDetailTransactions).Error
	if err != nil {
		return newDetailTransactions, err
	}
	return newDetailTransactions, nil
}
func (tr *transactionRepository) GetAllUserTransaction(userID int) ([]entity.Transaction, error) {
	var transaction []entity.Transaction
	tr.db.Where(userID).Find(&transaction)

	return transaction, nil
}
func (tr *transactionRepository) GetAllStoreTransaction(userID int) ([]entity.TransactionDetail, error) {
	var transactionDetail []entity.TransactionDetail
	err := tr.db.Table("transaction_details").Joins("join products on transaction_details.product_id = products.id").Where(
		"products.store_id = ?", userID,
	).Find(&transactionDetail).Error

	if err != nil || len(transactionDetail) == 0 {
		return transactionDetail, errors.New("not found")
	}

	return transactionDetail, nil
}
func (tr *transactionRepository) Callback(callback entity.Transaction) error {

	var callbackData entity.Transaction
	err := tr.db.Where("invoice_id = ?", callback.InvoiceID).Model(&callbackData).Updates(callback).Error

	if err != nil || callbackData.PaymentStatus != callback.PaymentStatus {
		return err
	}

	if callbackData.PaymentStatus == "EXPIRED" {

		var transaction entity.Transaction
		var transactionDetail []entity.TransactionDetail
		err = tr.db.Where("invoice_id = ?", callback.InvoiceID).First(&transaction).Error
		fmt.Println("ini transaction", transaction)
		if err != nil {
			return err
		}

		err = tr.db.Where("transaction_id = ?", transaction.ID).Find(&transactionDetail).Error

		if err != nil || len(transactionDetail) == 0 {
			return errors.New("error")
		}

		fmt.Println("ini detail", transactionDetail)

		for i := 0; i < len(transactionDetail); i++ {
			var product entity.Product
			err = tr.db.Where(" id = ?", transactionDetail[i].ProductID).First(&product).Error
			if err != nil {
				return err
			}

			if product.CategoryID != 1 {
				stock := product.Stock + transactionDetail[i].Quantity

				err = tr.db.Where("id = ?", product.ID).Model(&product).Update("stock", stock).Error

				if err != nil || product.Stock != stock {
					return err
				}
			}
		}
	}

	return nil
}

//transaction helper
func (tr *transactionRepository) GroomingStatusHelper(petID int, transactionDetailID uint) error {
	var groomingStatus entity.GroomingStatus
	groomingStatus = entity.GroomingStatus{
		PetID:               uint(petID),
		TransactionDetailID: transactionDetailID,
	}
	err := tr.db.Save(&groomingStatus).Error

	if err != nil {
		return errors.New("error")
	}

	return nil
}
func (tr *transactionRepository) PetValidator(petID int, userID int) error {
	var pet entity.Pet
	err := tr.db.Where("id = ? And user_id = ?", petID, userID).First(&pet).Error

	if err != nil {
		return err
	}

	return nil
}
func (tr *transactionRepository) GetProductByID(productID int) (entity.Product, error) {
	var product entity.Product
	err := tr.db.Where("id = ?", productID).First(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}
func (tr *transactionRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	var category entity.Category
	err := tr.db.Where("id = ?", categoryID).First(&category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}
func (tr *transactionRepository) GetUserByID(userID int) (entity.User, error) {
	var user entity.User
	err := tr.db.Where("id = ?", userID).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
func (tr *transactionRepository) UpdateStock(productID, stock int) error {
	var product entity.Product

	err := tr.db.Where("id = ?", productID).Model(&product).Update("stock", stock).Error

	if err != nil {
		return err
	}

	return nil
}
