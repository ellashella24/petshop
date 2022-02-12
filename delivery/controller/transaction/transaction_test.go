package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"petshop/constants"

	middlewares "petshop/delivery/middleware"

	"net/http"
	"net/http/httptest"
	"petshop/delivery/common"
	"petshop/entity"
	"testing"
	"time"
)

var jwt = ""

func TestSetup(t *testing.T) {
	userJWT, _ := middlewares.GenerateToken(1, "naufal@gmail.com", "admin")
	jwt = string(userJWT)

}
func TestCreate(t *testing.T) {
	t.Run(
		"1. Success Transaction Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				[]TransactionRequest{
					{
						ProductID: 1,
						Quantity:  5,
					},
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockTransactionRepository{})
			err := TransactionCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Success Transaction Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				[]TransactionRequest{
					{
						ProductID: 1,
						Quantity:  5,
						PetID:     1,
					},
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockGroomingTransactionRepository{})
			err := TransactionCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"3. Fail product not found Transaction Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				[]TransactionRequest{},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockFalseTransactionRepository{})
			err := TransactionCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
	t.Run(
		"4. Fail Product Transaction Bad Request Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				[]TransactionRequest{
					{
						ProductID: 1,
						Quantity:  1000,
					},
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockFalseTransactionRepository{})
			err := TransactionCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
	t.Run(
		"5. Fail Grooming Transaction Bad Request Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				[]TransactionRequest{
					{
						ProductID: 1,
						Quantity:  1,
						PetID:     0,
					},
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockGroomingTransactionRepository{})
			err := TransactionCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
	t.Run(
		"6. Fail Grooming Transaction Pet validator Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				[]TransactionRequest{
					{
						ProductID: 1,
						Quantity:  1,
						PetID:     6,
					},
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockGroomingFalseTransactionRepository{})
			err := TransactionCon.Create()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

}
func TestGetAllUserTransaction(t *testing.T) {
	t.Run(
		"1. Success Get All User Transaction Test", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockTransactionRepository{})
			err := TransactionCon.GetAllUserTransaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail Get All User Transaction Test", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockFalseTransactionRepository{})
			err := TransactionCon.GetAllUserTransaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
}
func TestGetAllStoreTransaction(t *testing.T) {
	t.Run(
		"1. Success Get All Store Transaction Test", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockTransactionRepository{})
			err := TransactionCon.GetAllStoreTransaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail Get All Store Transaction Test", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			res := httptest.NewRecorder()
			req.Header.Set("Bearer", jwt)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockFalseTransactionRepository{})
			err := TransactionCon.GetAllStoreTransaction()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseError

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)
}
func TestCallBack(t *testing.T) {
	t.Run(
		"1. Success Callback Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				entity.Transaction{
					InvoiceID:     "invoice",
					PaymentMethod: "BANK_TRANSFER",
					PaidAt:        time.Now(),
					PaymentStatus: "PAID",
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))

			res := httptest.NewRecorder()

			req.Header.Set("X-Callback-Token", constants.CallbackToken)
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockTransactionRepository{})
			err := TransactionCon.Callback()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
	t.Run(
		"2. Fail Callback Test", func(t *testing.T) {
			e := echo.New()

			requestBody, _ := json.Marshal(
				entity.Transaction{
					InvoiceID:     "invoice",
					PaymentMethod: "BANK_TRANSFER",
					PaidAt:        time.Now(),
					PaymentStatus: "PAID",
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/transaction")

			TransactionCon := NewTransactionController(mockTransactionRepository{})
			err := TransactionCon.Callback()(context)

			if err != nil {
				log.Error(err)
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Unauthorized", response.Message)
		},
	)
}

type mockTransactionRepository struct{}

func (mb mockTransactionRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
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
func (mb mockTransactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	return entity.TransactionDetail{
		ID:            1,
		TransactionID: 1,
		ProductID:     1,
		Quantity:      1,
	}, nil
}
func (mb mockTransactionRepository) GroomingStatusHelper(petID int, transactionDetailID uint) error {
	return nil
}
func (mb mockTransactionRepository) PetValidator(petID int, userID int) error {
	return nil
}
func (mb mockTransactionRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		ID:         1,
		Name:       "Whiskas",
		Price:      100000,
		Stock:      10,
		StoreID:    1,
		CategoryID: 2,
	}, nil
}
func (mb mockTransactionRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID:   2,
		Name: "Makanan",
	}, nil
}
func (mb mockTransactionRepository) GetUserByID(userID int) (entity.User, error) {
	return entity.User{
		ID:       1,
		Name:     "Naufal",
		Email:    "naufal@gmail.com",
		Password: "123",
		CityID:   1,
		Role:     "Admin",
	}, nil
}
func (mb mockTransactionRepository) Callback(callback entity.Transaction) error {
	return nil
}
func (mb mockTransactionRepository) UpdateStock(productID, stock int) error {
	return nil
}
func (mb mockTransactionRepository) GetAllUserTransaction(userID int) ([]entity.Transaction, error) {
	return []entity.Transaction{
		{
			ID:            1,
			UserID:        1,
			InvoiceID:     "invoice",
			PaymentMethod: "BANK_TRANSFER",
			PaidAt:        time.Now(),
			TotalPrice:    100000,
			PaymentStatus: "PENDING",
		},
	}, nil
}
func (mb mockTransactionRepository) GetAllStoreTransaction(userID int) ([]entity.TransactionDetail, error) {
	return []entity.TransactionDetail{
		{
			ID:            1,
			TransactionID: 1,
			ProductID:     1,
			Quantity:      10,
		},
	}, nil
}

type mockFalseTransactionRepository struct{}

func (mb mockFalseTransactionRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	return entity.Transaction{}, errors.New("Error")
}
func (mb mockFalseTransactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	return entity.TransactionDetail{}, errors.New("Error")
}
func (mb mockFalseTransactionRepository) GroomingStatusHelper(petID int, transactionDetailID uint) error {
	return errors.New("Error")
}
func (mb mockFalseTransactionRepository) PetValidator(petID int, userID int) error {
	return errors.New("Error")
}
func (mb mockFalseTransactionRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{}, errors.New("Error")
}
func (mb mockFalseTransactionRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{}, errors.New("Error")
}
func (mb mockFalseTransactionRepository) GetUserByID(userID int) (entity.User, error) {
	return entity.User{}, errors.New("Error")
}
func (mb mockFalseTransactionRepository) Callback(callback entity.Transaction) error {
	return errors.New("Error")
}
func (mb mockFalseTransactionRepository) UpdateStock(productID, stock int) error {
	return errors.New("Error")
}
func (mb mockFalseTransactionRepository) GetAllUserTransaction(userID int) ([]entity.Transaction, error) {
	return []entity.Transaction{}, errors.New("Error")
}
func (mb mockFalseTransactionRepository) GetAllStoreTransaction(userID int) ([]entity.TransactionDetail, error) {
	return []entity.TransactionDetail{}, errors.New("Error")
}

type mockGroomingTransactionRepository struct{}

func (mb mockGroomingTransactionRepository) Transaction(newTransactions entity.Transaction) (
	entity.Transaction, error,
) {
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
func (mb mockGroomingTransactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	return entity.TransactionDetail{
		ID:            1,
		TransactionID: 1,
		ProductID:     1,
		Quantity:      1,
	}, nil
}
func (mb mockGroomingTransactionRepository) GroomingStatusHelper(petID int, transactionDetailID uint) error {
	return nil
}
func (mb mockGroomingTransactionRepository) PetValidator(petID int, userID int) error {
	return nil
}
func (mb mockGroomingTransactionRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		ID:         1,
		Name:       "Whiskas",
		Price:      100000,
		Stock:      10,
		StoreID:    1,
		CategoryID: 1,
	}, nil
}
func (mb mockGroomingTransactionRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID:   1,
		Name: "Grooming",
	}, nil
}
func (mb mockGroomingTransactionRepository) GetUserByID(userID int) (entity.User, error) {
	return entity.User{
		ID:       1,
		Name:     "Naufal",
		Email:    "naufal@gmail.com",
		Password: "123",
		CityID:   1,
		Role:     "Admin",
	}, nil
}
func (mb mockGroomingTransactionRepository) Callback(callback entity.Transaction) error {
	return nil
}
func (mb mockGroomingTransactionRepository) UpdateStock(productID, stock int) error {
	return nil
}
func (mb mockGroomingTransactionRepository) GetAllUserTransaction(userID int) ([]entity.Transaction, error) {
	return []entity.Transaction{
		{
			ID:            1,
			UserID:        1,
			InvoiceID:     "invoice",
			PaymentMethod: "BANK_TRANSFER",
			PaidAt:        time.Now(),
			TotalPrice:    100000,
			PaymentStatus: "PENDING",
		},
	}, nil
}
func (mb mockGroomingTransactionRepository) GetAllStoreTransaction(userID int) ([]entity.TransactionDetail, error) {
	return []entity.TransactionDetail{
		{
			ID:            1,
			TransactionID: 1,
			ProductID:     1,
			Quantity:      10,
		},
	}, nil
}

type mockGroomingFalseTransactionRepository struct{}

func (mb mockGroomingFalseTransactionRepository) Transaction(newTransactions entity.Transaction) (
	entity.Transaction, error,
) {
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
func (mb mockGroomingFalseTransactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	return entity.TransactionDetail{
		ID:            1,
		TransactionID: 1,
		ProductID:     1,
		Quantity:      1,
	}, nil
}
func (mb mockGroomingFalseTransactionRepository) GroomingStatusHelper(petID int, transactionDetailID uint) error {
	return nil
}
func (mb mockGroomingFalseTransactionRepository) PetValidator(petID int, userID int) error {
	return errors.New("Error")
}
func (mb mockGroomingFalseTransactionRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		ID:         1,
		Name:       "Whiskas",
		Price:      100000,
		Stock:      10,
		StoreID:    1,
		CategoryID: 1,
	}, errors.New("error")
}
func (mb mockGroomingFalseTransactionRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID:   1,
		Name: "Grooming",
	}, nil
}
func (mb mockGroomingFalseTransactionRepository) GetUserByID(userID int) (entity.User, error) {
	return entity.User{
		ID:       1,
		Name:     "Naufal",
		Email:    "naufal@gmail.com",
		Password: "123",
		CityID:   1,
		Role:     "Admin",
	}, nil
}
func (mb mockGroomingFalseTransactionRepository) Callback(callback entity.Transaction) error {
	return nil
}
func (mb mockGroomingFalseTransactionRepository) UpdateStock(productID, stock int) error {
	return nil
}
func (mb mockGroomingFalseTransactionRepository) GetAllUserTransaction(userID int) ([]entity.Transaction, error) {
	return []entity.Transaction{
		{
			ID:            1,
			UserID:        1,
			InvoiceID:     "invoice",
			PaymentMethod: "BANK_TRANSFER",
			PaidAt:        time.Now(),
			TotalPrice:    100000,
			PaymentStatus: "PENDING",
		},
	}, nil
}
func (mb mockGroomingFalseTransactionRepository) GetAllStoreTransaction(userID int) (
	[]entity.TransactionDetail, error,
) {
	return []entity.TransactionDetail{
		{
			ID:            1,
			TransactionID: 1,
			ProductID:     1,
			Quantity:      10,
		},
	}, nil
}
