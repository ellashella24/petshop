package preference

import (
	"github.com/xendit/xendit-go"
)

var SendEmail = true

var SendNotifWith = xendit.InvoiceCustomerNotificationPreference{
	InvoicePaid: []string{"email"},
}
