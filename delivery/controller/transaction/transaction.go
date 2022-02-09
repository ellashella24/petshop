package transaction

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xendit/xendit-go"
	"net/http"
	"petshop/constants"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/helper"
	transRepo "petshop/repository/transaction"
	"strings"
	"time"
)

type TransactionController struct {
	transactionRepo transRepo.Transaction
}

func NewTransactionController(transactionRepo transRepo.Transaction) *TransactionController {
	return &TransactionController{transactionRepo}
}

func (tc TransactionController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var transactionRequest []TransactionRequest
		xendit.Opt.SecretKey = constants.XendToken

		fmt.Println(transactionRequest)
		// bind request data
		c.Bind(&transactionRequest)

		//get data from token
		userID := middleware.ExtractTokenUserID(c)

		//Create externalID

		year, month, day := time.Now().Date()
		hour, minute, second := time.Now().Clock()
		ExternalID := fmt.Sprint("Invoice-", userID, "-", year, month, day, hour, minute, second)

		ExternalID = strings.ReplaceAll(ExternalID, " ", "")
		//Invoice Item
		var invoiceItem = xendit.InvoiceItem{}
		var invoiceItems = []xendit.InvoiceItem{}
		var amount = 0

		for i := 0; i < len(transactionRequest); i++ {

			productID := transactionRequest[i].ProductID

			getProduct, _ := tc.transactionRepo.GetProductByID(productID)
			category, _ := tc.transactionRepo.GetCategoryByID(int(getProduct.CategoryID))

			//check Stock and update

			if category.Name != "Grooming" {
				if transactionRequest[i].Quantity > getProduct.Stock {
					return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
				} else {
					stock := getProduct.Stock - transactionRequest[i].Quantity
					tc.transactionRepo.UpdateStock(int(getProduct.ID), stock)
				}
			}
			//Check Pet
			if category.Name == "Grooming" && transactionRequest[i].PetID == 0 {

				return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())

			} else if category.Name == "Grooming" && transactionRequest[i].PetID != 0 {
				err := tc.transactionRepo.PetValidator(transactionRequest[i].PetID, userID)

				if err != nil {
					return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
				}
			}

			invoiceItem = xendit.InvoiceItem{
				Name:     getProduct.Name,
				Price:    float64(getProduct.Price),
				Quantity: transactionRequest[i].Quantity,
				Category: category.Name,
			}

			amount += getProduct.Price * transactionRequest[i].Quantity

			invoiceItems = append(invoiceItems, invoiceItem)
		}

		//xendit customer
		customerData, _ := tc.transactionRepo.GetUserByID(userID)

		customer := xendit.InvoiceCustomer{
			GivenNames: customerData.Name,
			Email:      customerData.Email,
		}

		res := helper.CreateInvoice(ExternalID, float64(amount), customerData.Email, invoiceItems, customer)
		//save to db
		transactionData := entity.Transaction{
			UserID:        uint(userID),
			InvoiceID:     ExternalID,
			PaymentURL:    res.InvoiceURL,
			TotalPrice:    int(res.Amount),
			PaymentStatus: res.Status,
		}

		transactionRes, err := tc.transactionRepo.Transaction(transactionData)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		for i := 0; i < len(transactionRequest); i++ {
			transactionDetailData := entity.TransactionDetail{
				TransactionID: transactionRes.ID,
				ProductID:     uint(transactionRequest[i].ProductID),
				Quantity:      transactionRequest[i].Quantity,
			}
			detail, err := tc.transactionRepo.TransactionDetail(transactionDetailData)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
			}

			if invoiceItems[i].Category == "Grooming" {
				tc.transactionRepo.GroomingStatusHelper(transactionRequest[i].PetID, detail.ID)
			} else {

			}
			{
			}

		}

		return c.JSON(http.StatusOK, common.SuccessResponse(transactionRes))
	}
}

func (tc TransactionController) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		callBackToken := headers.Get("X-Callback-Token")

		if callBackToken != constants.CallbackToken {
			return c.JSON(http.StatusUnauthorized, common.NewUnauthorized())
		}

		var callBackRequest CallbackRequest
		c.Bind(&callBackRequest)

		parsePaid, _ := time.Parse(time.RFC3339, callBackRequest.PaidAt)

		callBackData := entity.Transaction{
			InvoiceID:     callBackRequest.ExternalID,
			PaymentStatus: callBackRequest.Status,
			PaymentMethod: callBackRequest.PaymentMethod,
			PaidAt:        parsePaid,
		}

		tc.transactionRepo.Callback(callBackData)

		return c.JSON(http.StatusOK, common.SuccessResponse(callBackData))

	}
}
