package transaction

import (
	"errors"
	"fmt"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"gorm.io/gorm"
	"petshop/entity"

	"petshop/preference"
)

type Transaction interface {
	Transaction(newTransactions entity.Transaction) (entity.Transaction, error)
	TransactionDetail(newDetailTransactions entity.TransactionDetail) (entity.TransactionDetail, error)
	GroomingStatusHelper(petID int, transactionDetailID uint) error
	TransactionHelper(request []entity.TransactionDetail, userID int) (*xendit.Invoice, error)
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

func (tr *transactionRepository) Callback(callback entity.Transaction, invoiceID string) error {

	var callbackData entity.Transaction
	err := tr.db.Where("invoice_id = ?", invoiceID).Model(&callbackData).Update(
		"payment_status", callback.PaymentStatus,
	).Error

	if err != nil || callbackData.PaymentStatus != callback.PaymentStatus {
		return err
	}

	if callbackData.PaymentStatus == "EXPIRED" {
		var transactionDetail []entity.TransactionDetail

		err = tr.db.Where("transaction_id = ?", invoiceID).Find(&transactionDetail).Error

		if err != nil || len(transactionDetail) == 0 {
			return err
		}

		for i := 0; i < len(transactionDetail); i++ {
			var product entity.Product
			err = tr.db.Where(" id = ?", transactionDetail[i].ProductID).First(&product).Error
			if err != nil {
				return err
			}
			stock := product.Stock + transactionDetail[i].Quantity

			err = tr.db.Where("id = ?", product.ID).Model(&product).Update("stock", stock).Error

			if err != nil || product.Stock != stock {
				return err
			}

		}
	}

	return nil
}

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

func (tr *transactionRepository) TransactionHelper(
	request []entity.TransactionDetail, userID int,
) (*xendit.Invoice, error) {

	var invoiceData *xendit.Invoice

	/// Customer
	var customer xendit.InvoiceCustomer
	var user entity.User

	err := tr.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return invoiceData, err
	}

	customer.GivenNames = user.Name
	customer.Email = user.Email

	//Amount
	amount := 0
	//category
	category := ""

	//item
	var item = xendit.InvoiceItem{}
	var items = []xendit.InvoiceItem{}

	for i := 0; i < len(request); i++ {

		//Find product
		var product []entity.Product

		err = tr.db.Where("id = ?", request[i].ProductID).Find(&product).Error
		if err != nil && len(product) == 0 {
			return invoiceData, err
		}

		//Check Stock

		if request[i].Quantity > product[i].Stock {
			return invoiceData, err
		}

		//Find Product Category
		var categoryData entity.Category
		err = tr.db.Where("id = ?", product[i].CategoryID).First(&categoryData).Error
		if err != nil {
			return invoiceData, err
		}
		category = categoryData.Name

		//sub stock
		if category != "Grooming" {
			newStock := product[i].Stock - request[i].Quantity
			err = tr.db.Where("id = ?", request[i].ProductID).Model(&product).Update("stock", newStock).Error

			if err != nil {
				return invoiceData, err
			}
		}

		//count amount

		//item
		item = xendit.InvoiceItem{
			Name:     product[i].Name,
			Price:    float64(product[i].Price),
			Quantity: product[i].Price,
			Category: category,
		}
		items = append(items, item)
	}

	//InvoiceID generator

	data := invoice.CreateParams{
		ExternalID:                     ExternalID,
		Amount:                         float64(amount),
		PayerEmail:                     user.Email,
		CustomerNotificationPreference: preference.SendNotifWith,
		ShouldSendEmail:                &preference.SendEmail,
		Items:                          items,
		Customer:                       customer,
	}

	resp, err := invoice.Create(&data)

	if err != nil {
		fmt.Println(err)
	}

	return resp, nil

}
