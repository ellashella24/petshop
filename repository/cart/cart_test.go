package cart

import (
	"fmt"
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

func TestCreate(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)
	cartRepo := NewCartRepository(db)

	t.Run(
		"1. Succes case", func(t *testing.T) {

			mockRequest := entity.Cart{
				UserID:    1,
				ProductID: 1,
				Quantity:  3,
			}

			res, err := cartRepo.Create(mockRequest)

			assert.Equal(t, 3, res.Quantity)
			assert.Nil(t, err)

		},
	)
	t.Run(
		" 2.error case", func(t *testing.T) {
			mockRequest := entity.Cart{
				UserID:    1,
				ProductID: 2,
				Quantity:  1,
			}
			_, err := cartRepo.Create(mockRequest)

			assert.NotNil(t, err)

		},
	)
	t.Run(
		" 3.error case", func(t *testing.T) {
			mockRequest := entity.Cart{
				UserID:    100,
				ProductID: 1,
				Quantity:  1,
			}
			_, err := cartRepo.Create(mockRequest)
			fmt.Println(err)
			assert.NotNil(t, err)

		},
	)
}

func TestCheckCart(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)
	cartRepo := NewCartRepository(db)

	t.Run(
		"1. Succes case", func(t *testing.T) {

			mockRequest := entity.Cart{
				UserID:    1,
				ProductID: 1,
			}

			res, err := cartRepo.CheckCart(mockRequest)

			assert.Equal(t, 3, res.Quantity)
			assert.Nil(t, err)

		},
	)

	t.Run(
		" 2.error case", func(t *testing.T) {
			mockRequest := entity.Cart{
				UserID:    100,
				ProductID: 2,
			}
			_, err := cartRepo.CheckCart(mockRequest)

			assert.NotNil(t, err)

		},
	)
}

func TestGetAll(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)
	cartRepo := NewCartRepository(db)

	t.Run(
		"1. Succes case", func(t *testing.T) {

			mockRequest := 1

			res, _ := cartRepo.GetAll(mockRequest)

			assert.Equal(t, 3, res[0].Quantity)

		},
	)

	t.Run(
		" 2.error case", func(t *testing.T) {
			mockRequest := 10
			_, err := cartRepo.GetAll(mockRequest)

			assert.NotNil(t, err)

		},
	)
}

func TestUpdate(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)
	cartRepo := NewCartRepository(db)

	t.Run(
		"1. Succes case", func(t *testing.T) {

			mockRequest := entity.Cart{
				UserID:    1,
				ProductID: 1,
				Quantity:  10,
			}

			res, _ := cartRepo.Update(mockRequest)

			assert.Equal(t, 10, res.Quantity)

		},
	)

	t.Run(
		" 2.error case", func(t *testing.T) {
			mockRequest := entity.Cart{
				UserID:    100,
				ProductID: 100,
				Quantity:  100,
			}
			_, err := cartRepo.Update(mockRequest)

			assert.NotNil(t, err)

		},
	)
}

func TestTransaction(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

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
			res, err := cartRepo.Transaction(mockRequest)
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
			_, err := cartRepo.Transaction(mockRequest)
			assert.NotNil(t, err)
		},
	)
}

func TestTransactionDetail(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			mockRequest := entity.TransactionDetail{
				TransactionID:  1,
				ProductID:      1,
				Quantity:       5,
				GroomingStatus: entity.GroomingStatus{},
			}
			err := cartRepo.TransactionDetail(mockRequest)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			mockRequest := entity.TransactionDetail{
				ProductID:      1,
				Quantity:       5,
				GroomingStatus: entity.GroomingStatus{},
			}
			err := cartRepo.TransactionDetail(mockRequest)
			assert.NotNil(t, err)

		},
	)
}

func TestGetProductByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := cartRepo.GetProductByID(1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := cartRepo.GetProductByID(100)
			assert.NotNil(t, err)
		},
	)
}

func TestGetCategoryID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := cartRepo.GetCategoryByID(1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := cartRepo.GetCategoryByID(100)
			assert.NotNil(t, err)
		},
	)
}

func TestGetUserByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := cartRepo.GetUserByID(1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := cartRepo.GetUserByID(100)
			assert.NotNil(t, err)
		},
	)
}

func TestUpdateStock(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			err := cartRepo.UpdateStock(1, 10)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			err := cartRepo.UpdateStock(100, 10)
			assert.NotNil(t, err)
		},
	)
}

func TestDelete(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cartRepo := NewCartRepository(db)

	t.Run(
		"succes case", func(t *testing.T) {

			_, err := cartRepo.Delete(1, 1)
			assert.Nil(t, err)

		},
	)

	t.Run(
		"error case", func(t *testing.T) {

			_, err := cartRepo.Delete(100, 100)
			assert.NotNil(t, err)
		},
	)
}
