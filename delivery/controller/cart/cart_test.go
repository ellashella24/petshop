package cart

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"petshop/delivery/common"
	middlewares "petshop/delivery/middleware"
	"petshop/entity"
	"testing"
	"time"
)

var jwt = ""

func TestSetup(t *testing.T) {
	userJWT, _ := middlewares.GenerateToken(1, "naufal@gmail.com", "admin")
	jwt = string(userJWT)

}
func TestCartTansaction(t *testing.T) {
	t.Run(
		"1. Success Transaction Test(All Product)", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal([]CartTransactionRequest{})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockCartRepository{})
			err := CartCon.CartTansaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Success Transaction Test(Selected Product)", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal([]CartTransactionRequest{{ProductID: 1}})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockCartRepository{})
			err := CartCon.CartTansaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"3. Fail Transaction Test(All Product)", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal([]CartTransactionRequest{})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.CartTansaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
	t.Run(
		"4. Fail Transaction Test(Selected Product)", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal([]CartTransactionRequest{{ProductID: 1}})

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.CartTansaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
}
func TestCreate(t *testing.T) {
	t.Run(
		"1. Success Create Cart", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				CartRequest{
					UserID:    1,
					ProductID: 1,
					Quantity:  1,
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockCartRepository{})
			err := CartCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail internal server", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				CartRequest{
					UserID:    1,
					ProductID: 1,
					Quantity:  1,
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
	t.Run(
		"2. Fail internal server", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				CartRequest{
					UserID:    1,
					ProductID: 1,
					Quantity:  1,
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
}
func TestGetAll(t *testing.T) {
	t.Run(
		"1. Success Get All Cart", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockCartRepository{})
			err := CartCon.GetAll()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail Get All Cart", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.GetAll()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Not Found", response.Message)
		},
	)
}
func TestUpdate(t *testing.T) {
	t.Run(
		"1. Success Update Cart", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				CartRequest{
					ProductID: 1,
					Quantity:  10,
				},
			)

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockCartRepository{})
			err := CartCon.Update()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail Update Cart", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				CartRequest{
					ProductID: 100,
					Quantity:  10,
				},
			)

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.Update()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
}
func TestDelete(t *testing.T) {
	t.Run(
		"1. Success Delete Cart", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")
			context.SetParamNames("id")
			context.SetParamValues("1")
			CartCon := NewCartController(mockCartRepository{})
			err := CartCon.Delete()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail Update Cart Bad Request", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.Delete()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
	t.Run(
		"3. Fail Delete Cart Not Found", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")
			context.SetParamNames("id")
			context.SetParamValues("1000")
			CartCon := NewCartController(mockFalseCartRepository{})
			err := CartCon.Delete()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Not Found", response.Message)
		},
	)
}

//MOCK CART ALL PRODUCT SUCCES

type mockCartRepository struct{}

func (mc mockCartRepository) GetAll(userId int) ([]entity.Cart, error) {
	return []entity.Cart{
		{
			ID:        1,
			UserID:    1,
			ProductID: 1,
			Quantity:  1,
		},
	}, nil
}
func (mc mockCartRepository) Create(entity.Cart) (entity.Cart, error) {
	return entity.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  1,
	}, nil
}
func (mc mockCartRepository) Update(entity.Cart) (entity.Cart, error) {
	return entity.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
	}, nil
}
func (mc mockCartRepository) Delete(userId int, productId int) (entity.Cart, error) {
	return entity.Cart{}, nil
}
func (mc mockCartRepository) CheckCart(entity.Cart) (entity.Cart, error) {
	return entity.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  1,
	}, nil
}
func (mb mockCartRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	return entity.Transaction{
		ID:            1,
		UserID:        1,
		InvoiceID:     "invoice",
		PaymentMethod: "BANK_TRANSFER",
		PaymentURL:    "xendit.com",
		PaidAt:        time.Now(),
		TotalPrice:    100000,
		PaymentStatus: "PENDING",
	}, nil
}
func (mb mockCartRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) error {
	return nil
}
func (mb mockCartRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		ID:         1,
		Name:       "Whiskas",
		Price:      100000,
		Stock:      10,
		StoreID:    1,
		CategoryID: 2,
	}, nil
}
func (mb mockCartRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID:   2,
		Name: "Makanan",
	}, nil
}
func (mb mockCartRepository) GetUserByID(userID int) (entity.User, error) {
	return entity.User{
		ID:       1,
		Name:     "Naufal",
		Email:    "naufal@gmail.com",
		Password: "123",
		CityID:   1,
		Role:     "Admin",
	}, nil
}
func (mb mockCartRepository) UpdateStock(productID, stock int) error {
	return nil
}

//MOCK CART FALSE

type mockFalseCartRepository struct{}

func (mc mockFalseCartRepository) GetAll(userId int) ([]entity.Cart, error) {
	return []entity.Cart{
		{
			UserID:    1,
			ProductID: 1,
			Quantity:  100000,
		},
	}, errors.New("Error")
}
func (mc mockFalseCartRepository) Create(entity.Cart) (entity.Cart, error) {
	return entity.Cart{}, errors.New("Error")
}
func (mc mockFalseCartRepository) Update(entity.Cart) (entity.Cart, error) {
	return entity.Cart{}, errors.New("Error")
}
func (mc mockFalseCartRepository) Delete(userId int, productId int) (entity.Cart, error) {
	return entity.Cart{}, errors.New("Error")
}
func (mc mockFalseCartRepository) CheckCart(entity.Cart) (entity.Cart, error) {
	return entity.Cart{
		UserID:    1,
		ProductID: 1,
		Quantity:  100,
	}, errors.New("Error")
}
func (mb mockFalseCartRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	return entity.Transaction{}, errors.New("Error")
}
func (mb mockFalseCartRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) error {
	return errors.New("Error")
}
func (mb mockFalseCartRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		Name:       "Whiskas",
		Price:      50000,
		Stock:      8,
		StoreID:    1,
		CategoryID: 2,
	}, errors.New("Error")
}
func (mb mockFalseCartRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{}, errors.New("Error")
}
func (mb mockFalseCartRepository) GetUserByID(userID int) (entity.User, error) {
	return entity.User{}, errors.New("Error")
}
func (mb mockFalseCartRepository) UpdateStock(productID, stock int) error {
	return errors.New("Error")
}
