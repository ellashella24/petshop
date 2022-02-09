package helper

import (
	"fmt"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"gorm.io/gorm"
	"petshop/delivery/controller/transaction"
	"petshop/entity"
	"petshop/preference"
	"time"
)

var db *gorm.DB

func TransactionHelper(request []transaction.TransactionRequest, userID int) (*xendit.Invoice, error) {

	var invoiceData *xendit.Invoice

	/// Customer
	var customer xendit.InvoiceCustomer
	var user entity.User

	err := db.Where("id = ?", userID).First(&user).Error
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

		err = db.Where("id = ?", request[i].ProductID).Find(&product).Error
		if err != nil && len(product) == 0 {
			return invoiceData, err
		}

		//Check Stock

		if request[i].Quantity > product[i].Stock {
			return invoiceData, err
		}

		//Find Product Category
		var categoryData entity.Category
		err = db.Where("id = ?", product[i].CategoryID).First(&categoryData).Error
		if err != nil {
			return invoiceData, err
		}
		category = categoryData.Name

		//sub stock
		if category != "Grooming" {
			newStock := product[i].Stock - request[i].Quantity
			err = db.Where("id = ?", request[i].ProductID).Model(&product).Update("stock", newStock).Error

			if err != nil {
				return invoiceData, err
			}
		}

		//count amount
		amount += request[i].Quantity * product[i].Price

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
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	ExternalID := fmt.Sprint("Invoice-", userID, "-", year, month, day, hour, minute, second)

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
