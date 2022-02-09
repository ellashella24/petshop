package helper

import (
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
	"petshop/preference"
)

func CreateInvoice(
	invoiceID string, amount float64, email string, items []xendit.InvoiceItem, customer xendit.InvoiceCustomer,
) *xendit.Invoice {

	data := invoice.CreateParams{
		ExternalID:                     invoiceID,
		Amount:                         float64(amount),
		PayerEmail:                     email,
		CustomerNotificationPreference: preference.SendNotifWith,
		ShouldSendEmail:                &preference.SendEmail,
		Items:                          items,
		Customer:                       customer,
	}

	resp, _ := invoice.Create(&data)

	return resp

}
