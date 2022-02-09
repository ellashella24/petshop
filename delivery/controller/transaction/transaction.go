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
	transRepo "petshop/repository/transaction"
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

		// bind request data
		c.Bind(&transactionRequest)

		//get data from token
		userID := middleware.ExtractTokenUserID(c)

		//helper

		year, month, day := time.Now().Date()
		hour, minute, second := time.Now().Clock()
		ExternalID := fmt.Sprint("Invoice-", userID, "-", year, month, day, hour, minute, second)

		//save to db
		transactionData := entity.Transaction{
			UserID:        uint(userID),
			InvoiceID:     ExternalID,
			PaymentURL:    help.InvoiceURL,
			TotalPrice:    int(help.Amount),
			PaymentStatus: help.Status,
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
			detail, _ := tc.transactionRepo.TransactionDetail(transactionDetailData)

			if help.Items[i].Category == "Grooming" {
				err = tc.transactionRepo.GroomingStatusHelper(transactionRequest[i].PetID, detail.ID)
				if err != nil {
					return err
				}
			}

		}

		return c.JSON(http.StatusOK, common.SuccessResponse(help))
	}
}

//
//func (tc TransactionController) Callback() echo.HandlerFunc {
//	return func(c echo.Context) error {
//		req := c.Request()
//		headers := req.Header
//
//		callBackToken := headers.Get("X-Callback-Token")
//
//		if callBackToken != constants.CallbackToken {
//			return c.JSON(http.StatusUnauthorized, common.NewUnauthorized())
//		}
//
//		var callBackRequest CallbackRequest
//		c.Bind(&callBackRequest)
//
//		parsePaid, _ := time.Parse(time.RFC3339, callBackRequest.PaidAt)
//
//		callBackData := entities.Booking{
//			PaymentStatus: callBackRequest.Status,
//			PaymentMethod: callBackRequest.PaymentMethod,
//			PaidAt:        parsePaid,
//		}
//
//		callBack, _ := bc.bookRepo.Update(callBackRequest.ExternalID, callBackData)
//
//		return c.JSON(http.StatusOK, common.SuccessResponse(callBack))
//
//	}
//
//		return c.JSON(http.StatusOK, common.SuccessResponse(help))
//	}
