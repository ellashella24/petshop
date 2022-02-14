package cart

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xendit/xendit-go"
	"gorm.io/gorm"
	"net/http"
	"petshop/constants"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/helper"
	"petshop/repository/cart"
	"strconv"
	"strings"
	"time"
)

type CartController struct {
	CartRepo cart.Cart
}

func NewCartController(CartRepo cart.Cart) *CartController {
	return &CartController{CartRepo}
}

func (cc CartController) CartTansaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		var transactionRequest []CartTransactionRequest
		xendit.Opt.SecretKey = constants.XendToken

		// bind request data
		err := c.Bind(&transactionRequest)

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
		var cartData = entity.Cart{}
		var cartDatas = []entity.Cart{}
		var amount = 0

		if len(transactionRequest) == 0 {
			cart, _ := cc.CartRepo.GetAll(userID)

			for i := range cart {

				productID := cart[i].ProductID

				getProduct, _ := cc.CartRepo.GetProductByID(int(productID))
				category, _ := cc.CartRepo.GetCategoryByID(int(getProduct.CategoryID))

				//check Stock and update

				if cart[i].Quantity > getProduct.Stock {
					return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
				} else {
					stock := getProduct.Stock - cart[i].Quantity
					cc.CartRepo.UpdateStock(int(getProduct.ID), stock)
				}

				invoiceItem = xendit.InvoiceItem{
					Name:     getProduct.Name,
					Price:    float64(getProduct.Price),
					Quantity: cart[i].Quantity,
					Category: category.Name,
				}

				amount += getProduct.Price * cart[i].Quantity

				//delete cart
				cc.CartRepo.Delete(userID, int(productID))
				cartData = entity.Cart{
					Model:     gorm.Model{},
					ProductID: cart[i].ProductID,
					Quantity:  cart[i].Quantity,
				}

				cartDatas = append(cartDatas, cartData)
				invoiceItems = append(invoiceItems, invoiceItem)
			}
		} else {
			for i := 0; i < len(transactionRequest); i++ {
				var cartInput entity.Cart
				cartInput = entity.Cart{
					UserID:    uint(userID),
					ProductID: uint(transactionRequest[i].ProductID),
				}

				cartData, err = cc.CartRepo.CheckCart(cartInput)

				if err != nil {
					return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
				}

				productID := cartData.ProductID

				getProduct, _ := cc.CartRepo.GetProductByID(int(productID))
				category, _ := cc.CartRepo.GetCategoryByID(int(getProduct.CategoryID))

				//check Stock and update

				if cartData.Quantity > getProduct.Stock {
					return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
				} else {
					stock := getProduct.Stock - cartData.Quantity
					cc.CartRepo.UpdateStock(int(getProduct.ID), stock)
				}

				invoiceItem = xendit.InvoiceItem{
					Name:     getProduct.Name,
					Price:    float64(getProduct.Price),
					Quantity: cartData.Quantity,
					Category: category.Name,
				}

				amount += getProduct.Price * cartData.Quantity

				//delete cart
				cc.CartRepo.Delete(userID, int(productID))

				invoiceItems = append(invoiceItems, invoiceItem)
			}
		}

		//xendit customer
		customerData, _ := cc.CartRepo.GetUserByID(userID)

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

		transactionRes, err := cc.CartRepo.Transaction(transactionData)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		for i := 0; i < len(cartDatas); i++ {
			transactionDetailData := entity.TransactionDetail{
				TransactionID: transactionRes.ID,
				ProductID:     cartDatas[i].ProductID,
				Quantity:      cartDatas[i].Quantity,
			}
			err = cc.CartRepo.TransactionDetail(transactionDetailData)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
			}

		}

		return c.JSON(http.StatusOK, common.SuccessResponse(transactionRes))
	}
}

func (cc CartController) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var cartRequest CartRequest

		c.Bind(&cartRequest)

		userID := middleware.ExtractTokenUserID(c)

		cartData := entity.Cart{
			UserID:    uint(userID),
			ProductID: cartRequest.ProductID,
			Quantity:  cartRequest.Quantity,
		}

		data, err := cc.CartRepo.CheckCart(cartData)
		if err != nil {
			_, err = cc.CartRepo.Create(cartData)
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError, common.NewInternalServerErrorResponse(),
				)
			}
		} else {
			cartData.Quantity += data.Quantity
			_, err = cc.CartRepo.Update(cartData)
			if err != nil {
				return c.JSON(
					http.StatusInternalServerError, common.ErrorResponse(http.StatusInternalServerError, err.Error()),
				)
			}
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}

func (cc CartController) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		carts, err := cc.CartRepo.GetAll(userID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(carts))
	}

}

func (cc CartController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		var cartRequest CartRequest

		c.Bind(&cartRequest)

		data := entity.Cart{
			UserID:    uint(userID),
			ProductID: cartRequest.ProductID,
			Quantity:  cartRequest.Quantity,
		}

		_, err := cc.CartRepo.Update(data)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}

func (cc CartController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		productId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		_, err = cc.CartRepo.Delete(userID, productId)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}
