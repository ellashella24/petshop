package transaction

import (
	"fmt"
	"net/http"
	"petshop/constants"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/helper"
	transRepo "petshop/repository/transaction"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xendit/xendit-go"
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

		if len(invoiceItems) == 0 {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
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

		transactionRes, _ := tc.transactionRepo.Transaction(transactionData)

		for i := 0; i < len(transactionRequest); i++ {
			transactionDetailData := entity.TransactionDetail{
				TransactionID: transactionRes.ID,
				ProductID:     uint(transactionRequest[i].ProductID),
				Quantity:      transactionRequest[i].Quantity,
			}
			detail, _ := tc.transactionRepo.TransactionDetail(transactionDetailData)

			if invoiceItems[i].Category == "Grooming" {
				tc.transactionRepo.GroomingStatusHelper(transactionRequest[i].PetID, detail.ID)
			}

		}

		return c.JSON(http.StatusOK, common.SuccessResponse(transactionRes))
	}
}

func (tc TransactionController) GetAllUserTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		transaction, transactionDetail, err := tc.transactionRepo.GetAllUserTransaction(userID)
		if err != nil || len(transaction) == 0 {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		var response GetUserTransactionResponse
		var responses []GetUserTransactionResponse

		for i := 0; i < len(transaction); i++ {
			var transactionResponse = TransactionResponse{
				ID:            transaction[i].ID,
				UserID:        transaction[i].UserID,
				InvoiceID:     transaction[i].InvoiceID,
				PaymentMethod: transaction[i].PaymentMethod,
				PaymentURL:    transaction[i].PaymentURL,
				PaidAt:        transaction[i].PaidAt,
				TotalPrice:    transaction[i].TotalPrice,
				PaymentStatus: transaction[i].PaymentStatus,
			}

			var transactionDetailRes TransactionDetailResponse
			var transactionDetailResponses []TransactionDetailResponse

			for j := 0; j < len(transactionDetail[i]); j++ {
				transactionDetailRes = TransactionDetailResponse{
					TransactionID: transactionDetail[i][j].TransactionID,
					ProductID:     transactionDetail[i][j].ProductID,
					Quantity:      transactionDetail[i][j].Quantity,
				}
				transactionDetailResponses = append(transactionDetailResponses, transactionDetailRes)
			}

			response = GetUserTransactionResponse{
				Transaction:       transactionResponse,
				TransactionDetail: transactionDetailResponses,
			}

			responses = append(responses, response)

		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
	}
}

func (tc TransactionController) GetAllStoreTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)
		transactionDetail, transaction, err := tc.transactionRepo.GetAllStoreTransaction(userID)

		if err != nil || len(transactionDetail) == 0 {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		var respon GetStoreTransactionResponse
		var responses []GetStoreTransactionResponse
		counter := 0

		for i := 0; i < len(transaction); i++ {
			transactionRes := TransactionResponse{
				ID:            transaction[i].ID,
				UserID:        transaction[i].UserID,
				InvoiceID:     transaction[i].InvoiceID,
				PaymentMethod: transaction[i].PaymentMethod,
				PaymentURL:    transaction[i].PaymentURL,
				PaidAt:        transaction[i].PaidAt,
				TotalPrice:    transaction[i].TotalPrice,
				PaymentStatus: transaction[i].PaymentStatus,
			}

			transactionDetailRes := []TransactionDetailResponse{}

			for j := counter; j < len(transactionDetail); j++ {

				if transactionRes.ID == transactionDetail[j].TransactionID {
					transactionDetailData := TransactionDetailResponse{
						TransactionID: transactionDetail[j].TransactionID,
						ProductID:     transactionDetail[j].ProductID,
						Quantity:      transactionDetail[j].Quantity,
					}
					counter++
					transactionDetailRes = append(transactionDetailRes, transactionDetailData)
				} else {
					break
				}

			}

			respon = GetStoreTransactionResponse{
				Transaction:       transactionRes,
				TransactionDetail: transactionDetailRes,
			}

			responses = append(responses, respon)

		}
		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
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
		tc.transactionRepo.GroomingStatusHelperUpdate(callBackRequest.ExternalID)

		return c.JSON(http.StatusOK, common.SuccessResponse(callBackData))

	}
}
