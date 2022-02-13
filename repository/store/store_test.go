package store

import (
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.Migrator().DropTable(&entity.Store{})
	db.AutoMigrate(&entity.Store{})

	newCity := entity.City{}
	newCity.Name = "city1"

	db.Create(&newCity)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)
	newUser := entity.User{}
	newUser.Name = "user1"
	newUser.Email = "user1@mail.com"
	newUser.Password = string(hashedPassword)
	newUser.CityID = 1

	db.Create(&newUser)
}

func TestFindCityByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get city", func(t *testing.T) {
		res, err := storeRepo.FindCityByID(1)

		assert.Nil(t, err)
		assert.Equal(t, "city1", res.Name)
	})

	t.Run("2. Error get city", func(t *testing.T) {
		db.Migrator().DropTable(&entity.City{})
		res, _ := storeRepo.FindCityByID(100)

		assert.Equal(t, "", res.Name)
	})
}

func TestCreateStore(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	newCity := entity.City{}
	newCity.Name = "city2"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user2"), bcrypt.MinCost)
	newUser := entity.User{}
	newUser.Name = "user2"
	newUser.Email = "user2@mail.com"
	newUser.Password = string(hashedPassword)
	newUser.CityID = 1

	storeRepo := NewStoreRepository(db)

	t.Run("1. Success create store", func(t *testing.T) {
		db.Create(&newCity)
		db.Create(&newUser)
		newStore := entity.Store{}
		newStore.Name = "store2"
		newStore.CityID = 2
		newStore.UserID = 2

		res, err := storeRepo.CreateStore(newStore)

		assert.Nil(t, err)
		assert.Equal(t, newStore.Name, res.Name)
	})

	t.Run("2. Error create store", func(t *testing.T) {
		newStore := entity.Store{}
		newStore.Name = "store1"
		newStore.CityID = 0
		newStore.UserID = 0

		res, err := storeRepo.CreateStore(newStore)

		assert.NotNil(t, err)
		assert.Equal(t, 0, int(res.ID))
	})
}

func TestGetAllStoreByUserID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	storeRepo := NewStoreRepository(db)

	newStore := entity.Store{}
	newStore.Name = "store1"
	newStore.CityID = 1
	newStore.UserID = 1

	db.Create(&newStore)

	t.Run("1. Success get all store by user", func(t *testing.T) {
		res, err := storeRepo.GetAllStoreByUserID(1)

		assert.Nil(t, err)
		assert.Equal(t, "store1", res[0].Name)
	})

	t.Run("2. Error get store by user - user not found", func(t *testing.T) {
		res, err := storeRepo.GetAllStoreByUserID(10000)

		assert.NotNil(t, err)
		assert.Equal(t, true, len(res) == 0)
	})
}

func TestGetStoreProfile(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	newStore := entity.Store{}
	newStore.Name = "store1"
	newStore.CityID = 1
	newStore.UserID = 1

	db.Create(&newStore)

	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get store", func(t *testing.T) {
		res, err := storeRepo.GetStoreProfile(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "store1", res.Name)
	})

	t.Run("2. Error get store - store not found", func(t *testing.T) {
		res, err := storeRepo.GetStoreProfile(1000, 10000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestUpdateStoreProfile(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	newStore := entity.Store{}
	newStore.Name = " store1"
	newStore.UserID = 1

	db.Create(&newStore)

	storeRepo := NewStoreRepository(db)

	t.Run("1. Success update store", func(t *testing.T) {
		updatedStore := entity.Store{}
		updatedStore.Name = "store1 new"
		updatedStore.UserID = 1

		res, err := storeRepo.UpdateStoreProfile(1, 1, updatedStore)

		assert.Nil(t, err)
		assert.Equal(t, "store1 new", res.Name)
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
	config := config.GetConfig()
	db := util.InitDB(config)

	newStore := entity.Store{}
	newStore.Name = "store1"
	newStore.UserID = 1

	db.Create(&newStore)

	storeRepo := NewStoreRepository(db)

	t.Run("1. Success get store", func(t *testing.T) {
		_, err := storeRepo.DeleteStore(1, 1)

		assert.Nil(t, err)
	})

	t.Run("2. Error delete store - store not found", func(t *testing.T) {
		_, err := storeRepo.DeleteStore(1000, 10000)

		assert.NotNil(t, err)
	})
}
