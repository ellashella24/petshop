package transaction

import (
	"github.com/labstack/echo/v4"
	"github.com/xendit/xendit-go"
	"net/http"
	"petshop/constants"
	"petshop/delivery/common"
	"petshop/delivery/helper"
	"petshop/entity"
	"petshop/middlewares"
	"petshop/repository/transaction"
)

type TransactionController struct {
	transactionRepo transaction.Transaction
}

func NewTransactionController(transactionRepo transaction.Transaction) *TransactionController {
	return &TransactionController{transactionRepo}
}
func (tc TransactionController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var transactionRequest []TransactionRequest
		xendit.Opt.SecretKey = constants.XendToken

		// bind request data
		c.Bind(&transactionRequest)

		//get data from token
		userID := middlewares.NewAuth().ExtractTokenUserID(c)

		//helper
		help, err := helper.TransactionHelper(transactionRequest, userID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		//save to db
		transactionData := entity.Transaction{
			UserID:        uint(userID),
			InvoiceID:     help.ExternalID,
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
				helper.GroomingStatusHelper(transactionRequest[i].PetID, detail.ID)
			}

		}

		return c.JSON(http.StatusOK, common.SuccessResponse(help))
	}
}
