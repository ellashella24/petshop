package transaction

import (
	"errors"
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

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (tr *TransactionRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	err := tr.db.Save(&newTransactions).Error
	if err != nil {
		return newTransactions, err
	}
	return newTransactions, nil
}

func (tr *TransactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	err := tr.db.Save(&newDetailTransactions).Error
	if err != nil {
		return newDetailTransactions, err
	}
	return newDetailTransactions, nil
}
func (tr *TransactionRepository) GetAllUserTransaction(userID int) ([]entity.Transaction, error) {
	var transaction []entity.Transaction
	err := tr.db.Where(userID).Find(&transaction).Error

	if err != nil || len(transaction) == 0 {
		return transaction, err
	}
	return transaction, nil
}
func (tr *TransactionRepository) GetAllStoreTransaction(userID int) ([]entity.TransactionDetail, error) {
	var transactionDetail []entity.TransactionDetail
	err := tr.db.Table("transaction_details").Joins("join products on transaction_details.product_id = products.id").Where(
		"products.store_id = ?", userID,
	).Find(&transactionDetail).Error

	if err != nil || len(transactionDetail) == 0 {
		return transactionDetail, errors.New("not found")
	}

	return transactionDetail, nil
}
func (tr *TransactionRepository) Callback(callback entity.Transaction) error {

	var callbackData entity.Transaction
	err := tr.db.Where("invoice_id = ?", callback.InvoiceID).Model(&callbackData).Updates(callback).Error

	if err != nil || callbackData.InvoiceID == "" {
		return errors.New("Error")
	}

	if callbackData.PaymentStatus == "EXPIRED" {

		var transaction entity.Transaction
		var transactionDetail []entity.TransactionDetail
		err = tr.db.Where("invoice_id = ?", callback.InvoiceID).First(&transaction).Error

		if err != nil {
			return err
		}

		err = tr.db.Where("transaction_id = ?", transaction.ID).Find(&transactionDetail).Error

		if err != nil || len(transactionDetail) == 0 {
			return errors.New("error")
		}

		for i := 0; i < len(transactionDetail); i++ {
			var product entity.Product
			tr.db.Where(" id = ?", transactionDetail[i].ProductID).First(&product)

			if product.CategoryID != 1 {
				stock := product.Stock + transactionDetail[i].Quantity

				tr.db.Where("id = ?", product.ID).Model(&product).Update("stock", stock)

			}
		}
	}

	return nil
}

//transaction helper
func (tr *TransactionRepository) GroomingStatusHelper(petID int, transactionDetailID uint) error {
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
func (tr *TransactionRepository) PetValidator(petID int, userID int) error {
	var pet entity.Pet
	err := tr.db.Where("id = ? And user_id = ?", petID, userID).First(&pet).Error

	if err != nil {
		return err
	}

	return nil
}
func (tr *TransactionRepository) GetProductByID(productID int) (entity.Product, error) {
	var product entity.Product
	err := tr.db.Where("id = ?", productID).First(&product).Error

	if err != nil {
		return product, err
	}

	return product, nil
}
func (tr *TransactionRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	var category entity.Category
	err := tr.db.Where("id = ?", categoryID).First(&category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}
func (tr *TransactionRepository) GetUserByID(userID int) (entity.User, error) {
	var user entity.User
	err := tr.db.Where("id = ?", userID).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
func (tr *TransactionRepository) UpdateStock(productID, stock int) error {
	var product entity.Product

	err := tr.db.Where("id = ?", productID).First(&product).Error

	if err != nil {
		return err
	}

	tr.db.Model(&product).Update("stock", stock)

	return nil
}
