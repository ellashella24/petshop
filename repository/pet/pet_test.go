package pet

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

func TestCreatePet(t *testing.T) {
	petRepo := NewPetRepository(db)

	t.Run("1. Success create pet", func(t *testing.T) {
		newPet := entity.Pet{}
		newPet.Name = "Pet 1"
		newPet.UserID = 1

		res, err := petRepo.CreatePet(newPet)

		assert.Nil(t, err)
		assert.Equal(t, newPet.Name, res.Name)
	})

	t.Run("2. Error create pet", func(t *testing.T) {
		newPet := entity.Pet{}
		newPet.Name = "Pet 1"
		newPet.UserID = 0

		res, err := petRepo.CreatePet(newPet)

		assert.NotNil(t, err)
		assert.Equal(t, 0, int(res.ID))
	})
}

func TestGetAllPetByUserID(t *testing.T) {
	petRepo := NewPetRepository(db)

	t.Run("1. Success get all pet by user", func(t *testing.T) {
		res, err := petRepo.GetAllPetByUserID(1)

		assert.Nil(t, err)
		assert.Equal(t, "Pet 1", res[0].Name)
	})

	t.Run("2. Error get pet by user - user not found", func(t *testing.T) {
		res, err := petRepo.GetAllPetByUserID(10000)

		assert.NotNil(t, err)
		assert.Equal(t, true, len(res) == 0)
	})
}

func TestGetPetProfileByID(t *testing.T) {
	petRepo := NewPetRepository(db)

	t.Run("1. Success get pet", func(t *testing.T) {
		res, err := petRepo.GetPetProfileByID(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "Pet 1", res.Name)
	})

	t.Run("2. Error get pet - pet not found", func(t *testing.T) {
		res, err := petRepo.GetPetProfileByID(1000, 10000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestGetGroomingStatusByPetID(t *testing.T) {
	newStore := entity.Store{}
	newStore.Name = "Store 1"
	newStore.CityID = 1
	newStore.UserID = 1

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

	gs1 := entity.GroomingStatus{}
	gs1.PetID = 1
	gs1.Status = "TELAH DIBAYAR"
	gs1.TransactionDetailID = 1

	db.Create(&newStore)
	db.Create(&product1)
	db.Create(&transaction1)
	db.Create(&transactionDetail1)
	db.Create(&gs1)

	petRepo := NewPetRepository(db)

	t.Run("1. Success get grooming status", func(t *testing.T) {
		res, err := petRepo.GetGroomingStatusByPetID(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "TELAH DIBAYAR", res.Status)
	})

	t.Run("2. Error get grooming status", func(t *testing.T) {
		res, err := petRepo.GetGroomingStatusByPetID(1000, 1000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Status)
	})
}

func TestUpdateFinalGroomingStatus(t *testing.T) {
	petRepo := NewPetRepository(db)

	t.Run("1. Error get grooming status - pet not found", func(t *testing.T) {
		res, err := petRepo.UpdateFinalGroomingStatus(1000, 1000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Status)
	})

	t.Run("2. Error get grooming status - grooming status not 'DELIVERY KE USER'", func(t *testing.T) {
		res, err := petRepo.UpdateFinalGroomingStatus(1, 1)

		assert.NotNil(t, err)
		assert.Equal(t, "TELAH DIBAYAR", res.Status)
	})

	t.Run("3. Success get grooming status - change status to 'SELESAI'", func(t *testing.T) {
		gs1 := entity.GroomingStatus{}

		db.Where("pet_id = 1").Find(&gs1)

		db.Model(&gs1).Update("status", "DELIVERY KE USER")

		res, err := petRepo.UpdateFinalGroomingStatus(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "SELESAI", res.Status)
	})
}

func TestUpdatePetProfile(t *testing.T) {
	petRepo := NewPetRepository(db)

	t.Run("1. Success update pet", func(t *testing.T) {
		updatedPet := entity.Pet{}
		updatedPet.Name = "Pet 1 new"
		updatedPet.UserID = 1

		res, err := petRepo.UpdatePetProfile(1, 1, updatedPet)

		assert.Nil(t, err)
		assert.Equal(t, "Pet 1 new", res.Name)
	})

	t.Run("2. Error update pet - pet not found", func(t *testing.T) {
		updatedPet := entity.Pet{}
		updatedPet.Name = "Pet 1 new"
		updatedPet.UserID = 1

		res, err := petRepo.UpdatePetProfile(1000, 1000, updatedPet)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestDeletePet(t *testing.T) {
	petRepo := NewPetRepository(db)

	t.Run("1. Success get pet", func(t *testing.T) {
		_, err := petRepo.DeletePet(1, 1)

		assert.Nil(t, err)
	})

	t.Run("2. Error delete pet - pet not found", func(t *testing.T) {
		_, err := petRepo.DeletePet(1000, 10000)

		assert.NotNil(t, err)
	})
}
