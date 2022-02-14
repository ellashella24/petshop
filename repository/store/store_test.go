package store

import (
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db = util.InitDB(config)

	db.Migrator().DropTable(&entity.TransactionDetail{})
	db.Migrator().DropTable(&entity.Transaction{})
	db.Migrator().DropTable(&entity.GroomingStatus{})
	db.Migrator().DropTable(&entity.StockHistory{})
	db.Migrator().DropTable(&entity.Cart{})
	db.Migrator().DropTable(&entity.Product{})
	db.Migrator().DropTable(&entity.Category{})
	db.Migrator().DropTable(&entity.Store{})
	db.Migrator().DropTable(&entity.Pet{})
	db.Migrator().DropTable(&entity.User{})
	db.Migrator().DropTable(&entity.City{})

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

	newCity := entity.City{
		Name: "City 1",
	}

	db.Create(&newCity)

	newCategory := entity.Category{
		Name: "Grooming",
	}

	db.Create(&newCategory)

	newCategory1 := entity.Category{
		Name: "Makanan",
	}

	db.Create(&newCategory1)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)

	newUser := entity.User{
		Name:     "User 1",
		Email:    "user1@mail.com",
		Password: string(hashedPassword),
		CityID:   1,
	}

	db.Create(&newUser)
}

func TestFindCityByID(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get city", func(t *testing.T) {
		res, err := storeRepo.FindCityByID(1)

		assert.Nil(t, err)
		assert.Equal(t, "City 1", res.Name)
	})

	t.Run("2. Error get city", func(t *testing.T) {
		res, _ := storeRepo.FindCityByID(100)

		assert.Equal(t, "", res.Name)
	})
}

func TestCreateStore(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success create store", func(t *testing.T) {
		newStore := entity.Store{}
		newStore.Name = "Store 1"
		newStore.CityID = 1
		newStore.UserID = 1

		res, err := storeRepo.CreateStore(newStore)

		assert.Nil(t, err)
		assert.Equal(t, newStore.Name, res.Name)
	})

	t.Run("2. Error create store", func(t *testing.T) {
		newStore := entity.Store{}
		newStore.Name = "store0"
		newStore.CityID = 0
		newStore.UserID = 0

		res, err := storeRepo.CreateStore(newStore)

		assert.NotNil(t, err)
		assert.Equal(t, 0, int(res.ID))
	})
}

func TestGetAllStoreByUserID(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get all store by user", func(t *testing.T) {
		res, err := storeRepo.GetAllStoreByUserID(1)

		assert.Nil(t, err)
		assert.Equal(t, "Store 1", res[0].Name)
	})

	t.Run("2. Error get store by user - user not found", func(t *testing.T) {
		res, err := storeRepo.GetAllStoreByUserID(10000)

		assert.NotNil(t, err)
		assert.Equal(t, true, len(res) == 0)
	})
}

func TestGetStoreProfile(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get store", func(t *testing.T) {
		res, err := storeRepo.GetStoreProfile(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "Store 1", res.Name)
	})

	t.Run("2. Error get store - store not found", func(t *testing.T) {
		res, err := storeRepo.GetStoreProfile(1000, 10000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestGetListTransactionByStoreID(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	pet1 := entity.Pet{}
	pet1.Name = "Pet 1"
	pet1.UserID = 1

	product1 := entity.Product{}
	product1.Name = "Grooming"
	product1.Price = 100000
	product1.Stock = 1
	product1.CategoryID = 1
	product1.StoreID = 1

	transaction1 := entity.Transaction{}
	transaction1.UserID = uint(1)
	transaction1.InvoiceID = "Invoice-1"
	transaction1.PaymentMethod = "Bank Transfer"
	transaction1.PaymentURL = "payment-url-1"
	transaction1.PaidAt = time.Now()
	transaction1.TotalPrice = 800000
	transaction1.PaymentStatus = "Paid"

	transactionDetail1 := entity.TransactionDetail{}
	transactionDetail1.TransactionID = 1
	transactionDetail1.ProductID = 1
	transactionDetail1.Quantity = 1

	db.Create(&pet1)
	db.Create(&product1)
	db.Create(&transaction1)
	db.Create(&transactionDetail1)

	t.Run("1. Success get list transaction", func(t *testing.T) {
		transcation, transactionDetail, product, err := storeRepo.GetListTransactionByStoreID(1)

		assert.Nil(t, err)
		assert.Equal(t, "Invoice-1", transcation[0].InvoiceID)
		assert.Equal(t, 1, int(transactionDetail[0].ProductID))
		assert.Equal(t, 1, int(product[0].ID))
	})

	t.Run("2. Error get list transaction", func(t *testing.T) {
		transcation, transactionDetail, product, err := storeRepo.GetListTransactionByStoreID(1000)

		assert.NotNil(t, err)
		assert.Equal(t, []entity.Transaction([]entity.Transaction(nil)), transcation)
		assert.Equal(t, []entity.TransactionDetail([]entity.TransactionDetail{}), transactionDetail)
		assert.Equal(t, []entity.Product([]entity.Product(nil)), product)
	})
}

func TestGetGroomingStatusByPetID(t *testing.T) {
	gs1 := entity.GroomingStatus{}
	gs1.PetID = 1
	gs1.Status = "TELAH DIBAYAR"
	gs1.TransactionDetailID = 1

	db.Create(&gs1)

	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get grooming status", func(t *testing.T) {
		res, err := storeRepo.GetGroomingStatusByPetID(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "TELAH DIBAYAR", res.Status)
	})

	t.Run("2. Error get grooming status", func(t *testing.T) {
		res, err := storeRepo.GetGroomingStatusByPetID(1000, 1000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Status)
	})
}

func TestUpdateGroomingStatusByPetID(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success update grooming status - change status to 'PROSES PICKUP'", func(t *testing.T) {
		res, err := storeRepo.UpdateGroomingStatus(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "PROSES PICKUP", res.Status)
	})

	t.Run("2. Success update grooming status - change status to 'PROSES GROOMING'", func(t *testing.T) {
		res, err := storeRepo.UpdateGroomingStatus(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "PROSES GROOMING", res.Status)
	})

	t.Run("3. Success update grooming status - change status to 'DELIVERY KE USER'", func(t *testing.T) {
		res, err := storeRepo.UpdateGroomingStatus(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "DELIVERY KE USER", res.Status)
	})

	t.Run("4. Error update grooming status - status has already 'DELIVERY KE USER'", func(t *testing.T) {
		res, err := storeRepo.UpdateGroomingStatus(1, 1)

		assert.NotNil(t, err)
		assert.Equal(t, "DELIVERY KE USER", res.Status)
	})

	t.Run("5. Error get grooming status - pet & store not found", func(t *testing.T) {
		res, err := storeRepo.UpdateGroomingStatus(1000, 1000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Status)
	})
}

func TestUpdateStoreProfile(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success update store", func(t *testing.T) {
		updatedStore := entity.Store{}
		updatedStore.Name = "Store 1 new"
		updatedStore.UserID = 1

		res, err := storeRepo.UpdateStoreProfile(1, 1, updatedStore)

		assert.Nil(t, err)
		assert.Equal(t, "Store 1 new", res.Name)
	})

	t.Run("2. Error update store - store not found", func(t *testing.T) {
		updatedStore := entity.Store{}
		updatedStore.Name = " store1 new"
		updatedStore.UserID = 1

		res, err := storeRepo.UpdateStoreProfile(1000, 1000, updatedStore)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestDeleteStore(t *testing.T) {
	storeRepo := NewStoreRepository(db)

	t.Run("1. Success delete store", func(t *testing.T) {
		_, err := storeRepo.DeleteStore(1, 1)

		assert.Nil(t, err)
	})

	t.Run("2. Error delete store - store not found", func(t *testing.T) {
		_, err := storeRepo.DeleteStore(1000, 10000)

		assert.NotNil(t, err)
	})
}
