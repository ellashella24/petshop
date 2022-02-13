package transaction

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)
	db.Migrator().DropTable(&entity.TransactionDetail{})
	db.Migrator().DropTable(&entity.Transaction{})
	db.Migrator().DropTable(&entity.GroomingStatus{})
	db.Migrator().DropTable(&entity.StockHistory{})
	db.Migrator().DropTable(&entity.Cart{})
	db.Migrator().DropTable(&entity.Product{})
	db.Migrator().DropTable(&entity.Category{})
	db.Migrator().DropTable(&entity.Store{})
	db.Migrator().DropTable(&entity.User{})
	db.Migrator().DropTable(&entity.City{})
	db.Migrator().DropTable(&entity.Pet{})

	db.AutoMigrate(&entity.City{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Pet{})
	db.AutoMigrate(&entity.Store{})
	db.AutoMigrate(&entity.Category{})
	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.Cart{})
	db.AutoMigrate(&entity.StockHistory{})
	db.AutoMigrate(&entity.GroomingStatus{})
	db.AutoMigrate(&entity.Transaction{})
	db.AutoMigrate(&entity.TransactionDetail{})

	city := entity.City{}
	city.Name = "city1"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)
	admin := entity.User{
		Name:        "admin1",
		Email:       "fivonej388@kruay.com",
		Password:    string(hashedPassword),
		CityID:      1,
		Role:        "admin",
		Transaction: nil,
		Cart:        nil,
	}

	store := entity.Store{}
	store.Name = "store1"
	store.UserID = 1
	store.CityID = 1

	category := entity.Category{
		Name:    "Grooming",
		Product: nil,
	}
	category1 := entity.Category{
		Name:    "Makanan",
		Product: nil,
	}
	pet := entity.Pet{
		Name:   "Kucing",
		UserID: 1,
	}

	product := entity.Product{
		Name:       "Whiskas",
		Price:      100000,
		Stock:      100,
		StoreID:    1,
		CategoryID: 2,
	}
	product1 := entity.Product{
		Name:       "Shampoan",
		Price:      50000,
		Stock:      1,
		StoreID:    1,
		CategoryID: 1,
	}

	db.Create(&city)
	db.Create(&admin)
	db.Create(&store)
	db.Create(&category)
	db.Create(&category1)
	db.Create(&pet)
	db.Create(&product)
	db.Create(&product1)

}

func TestTransaction(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := entity.Transaction{
				UserID:            1,
				InvoiceID:         "invoice",
				PaymentMethod:     "BANK_TRANSFER",
				PaidAt:            time.Now(),
				TotalPrice:        150000,
				PaymentStatus:     "PENDING",
				TransactionDetail: nil,
			}
			res, err := transactionRepo.Transaction(mockRequest)
			assert.Nil(t, err)
			assert.Equal(t, uint(1), res.UserID)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			mockRequest := entity.Transaction{
				InvoiceID:         "invoice",
				PaymentMethod:     "BANK_TRANSFER",
				PaidAt:            time.Now(),
				TotalPrice:        150000,
				PaymentStatus:     "PENDING",
				TransactionDetail: nil,
			}
			_, err := transactionRepo.Transaction(mockRequest)
			assert.NotNil(t, err)
		},
	)
}

func TestTransactionDetail(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := entity.TransactionDetail{
				TransactionID:  1,
				ProductID:      1,
				Quantity:       5,
				GroomingStatus: entity.GroomingStatus{},
			}
			res, err := transactionRepo.TransactionDetail(mockRequest)
			assert.Nil(t, err)
			assert.Equal(t, uint(1), res.ProductID)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			mockRequest := entity.TransactionDetail{
				ProductID:      1,
				Quantity:       5,
				GroomingStatus: entity.GroomingStatus{},
			}
			_, err := transactionRepo.TransactionDetail(mockRequest)
			assert.NotNil(t, err)

		},
	)
}

func TestGetAllUserTransaction(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := 1
			transaction, _, err := transactionRepo.GetAllUserTransaction(mockRequest)
			assert.Nil(t, err)
			assert.Equal(t, uint(1), transaction[0].UserID)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			mockRequest := 100
			transaction, _, _ := transactionRepo.GetAllUserTransaction(mockRequest)
			assert.Equal(t, 0, len(transaction))
		},
	)
}

func TestGetAllStoreTransaction(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	var transaction = entity.Transaction{
		ID:                2,
		UserID:            1,
		InvoiceID:         "invoice12",
		PaymentMethod:     "",
		PaymentURL:        "",
		PaidAt:            time.Time{},
		TotalPrice:        100000,
		PaymentStatus:     "Pending",
		TransactionDetail: nil,
	}
	var transactionDetail = entity.TransactionDetail{
		TransactionID:  2,
		ProductID:      1,
		Quantity:       1,
		GroomingStatus: entity.GroomingStatus{},
	}

	db.Save(&transaction)
	db.Save(&transactionDetail)

	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := 1
			transactionDetail, transaction, err := transactionRepo.GetAllStoreTransaction(mockRequest)
			assert.Nil(t, err)
			assert.Equal(t, 5, transactionDetail[0].Quantity)
			assert.Equal(t, uint(1), transaction[0].ID)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			mockRequest := 100
			transactionDetail, _, _ := transactionRepo.GetAllStoreTransaction(mockRequest)
			assert.Equal(t, 0, len(transactionDetail))
		},
	)
}

func TestCallback(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := entity.Transaction{
				UserID:        1,
				InvoiceID:     "invoice",
				PaymentMethod: "BANK_TRANSFER",
				PaidAt:        time.Now(),
				TotalPrice:    150000,
				PaymentStatus: "EXPIRED",
			}
			err := transactionRepo.Callback(mockRequest)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case 1", func(t *testing.T) {

			mockRequest := entity.Transaction{

				PaymentMethod: "BANK_TRANSFER",
				PaidAt:        time.Now(),
				TotalPrice:    150000,
				PaymentStatus: "EXPIRED",
			}
			err := transactionRepo.Callback(mockRequest)
			assert.NotNil(t, err)
		},
	)
	t.Run(
		"error case 2", func(t *testing.T) {

			mockRequest := entity.Transaction{
				UserID:        1,
				InvoiceID:     "sgeghetheh",
				PaymentMethod: "BANK_TRANSFER",
				PaidAt:        time.Now(),
				TotalPrice:    150000,
				PaymentStatus: "EXPIRED",
			}
			err := transactionRepo.Callback(mockRequest)
			assert.NotNil(t, err)
		},
	)

	t.Run(
		"error case 3", func(t *testing.T) {

			var transactionDetail entity.TransactionDetail
			mockRequest := entity.Transaction{
				UserID:        1,
				InvoiceID:     "invoice",
				PaymentMethod: "BANK_TRANSFER",
				PaidAt:        time.Now(),
				TotalPrice:    150000,
				PaymentStatus: "EXPIRED",
			}

			db.Where("id", 1).Delete(&transactionDetail)

			err := transactionRepo.Callback(mockRequest)
			assert.NotNil(t, err)
		},
	)

}

func TestGroomingStatusHelper(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			var transactionDetail entity.TransactionDetail
			db.Where("id = ?", 1).Model(&transactionDetail).Update("deleted_at", "")

			mockRequest := entity.GroomingStatus{
				PetID:               1,
				TransactionDetailID: 1,
			}
			err := transactionRepo.GroomingStatusHelper(int(mockRequest.PetID), mockRequest.TransactionDetailID)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			err := transactionRepo.GroomingStatusHelper(100, 100)
			assert.NotNil(t, err)
		},
	)
}

func TestGroomingStatusHelperUpdate(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	mockRequest := entity.Transaction{
		UserID:            1,
		InvoiceID:         "invoice1",
		PaymentMethod:     "BANK_TRANSFER",
		PaidAt:            time.Now(),
		TotalPrice:        50000,
		PaymentStatus:     "PENDING",
		TransactionDetail: nil,
	}

	mockRequest1 := entity.TransactionDetail{
		TransactionID: 3,
		ProductID:     2,
	}

	transactionRepo.Transaction(mockRequest)
	transactionRepo.TransactionDetail(mockRequest1)
	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := "invoice1"
			err := transactionRepo.GroomingStatusHelperUpdate(mockRequest)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			mockRequest := "ghetdh"
			err := transactionRepo.GroomingStatusHelperUpdate(mockRequest)
			assert.NotNil(t, err)

		},
	)
}

func TestPetValidator(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			err := transactionRepo.PetValidator(1, 1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			err := transactionRepo.PetValidator(100, 100)
			assert.NotNil(t, err)
		},
	)
}

func TestGetProductByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := transactionRepo.GetProductByID(1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := transactionRepo.GetProductByID(100)
			assert.NotNil(t, err)
		},
	)
}

func TestGetCategoryID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := transactionRepo.GetCategoryByID(1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := transactionRepo.GetCategoryByID(100)
			assert.NotNil(t, err)
		},
	)
}

func TestGetUserByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := transactionRepo.GetUserByID(1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := transactionRepo.GetUserByID(100)
			assert.NotNil(t, err)
		},
	)
}

func TestUpdateStock(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	transactionRepo := NewTransactionRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			err := transactionRepo.UpdateStock(1, 10)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			err := transactionRepo.UpdateStock(100, 10)
			assert.NotNil(t, err)
		},
	)
}
