package pet

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

	db.Migrator().DropTable(&entity.Pet{})
	db.AutoMigrate(&entity.Pet{})

	newCity := entity.City{}
	newCity.Name = "city1"

	db.Create(&newCity)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)
	newUser := entity.User{}
	newUser.Name = "user1"
	newUser.Email = "user1@mail.com"
	newUser.Password = string(hashedPassword)
	newUser.CityID = 1
}

func TestCreatePet(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	petRepo := NewPetRepository(db)

	t.Run("1. Success create pet", func(t *testing.T) {
		newPet := entity.Pet{}
		newPet.Name = "pet1"
		newPet.UserID = 1

		res, err := petRepo.CreatePet(newPet)

		assert.Nil(t, err)
		assert.Equal(t, newPet.Name, res.Name)
	})

	t.Run("2. Error create pet", func(t *testing.T) {
		newPet := entity.Pet{}
		newPet.Name = "pet1"
		newPet.UserID = 0

		res, err := petRepo.CreatePet(newPet)

		assert.NotNil(t, err)
		assert.Equal(t, 0, int(res.ID))
	})
}

func TestGetAllPetByUserID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	petRepo := NewPetRepository(db)

	newPet := entity.Pet{}
	newPet.Name = "pet1"
	newPet.UserID = 1

	db.Create(&newPet)

	t.Run("1. Success get all pet by user", func(t *testing.T) {
		res, err := petRepo.GetAllPetByUserID(1)

		assert.Nil(t, err)
		assert.Equal(t, "pet1", res[0].Name)
	})

	t.Run("2. Error get pet by user - user not found", func(t *testing.T) {
		res, err := petRepo.GetAllPetByUserID(10000)

		assert.NotNil(t, err)
		assert.Equal(t, true, len(res) == 0)
	})
}

func TestGetPetProfileByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	newPet := entity.Pet{}
	newPet.Name = "pet1"
	newPet.UserID = 1

	db.Create(&newPet)

	petRepo := NewPetRepository(db)

	t.Run("1. Success get pet", func(t *testing.T) {
		res, err := petRepo.GetPetProfileByID(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, "pet1", res.Name)
	})

	t.Run("2. Error get pet - pet not found", func(t *testing.T) {
		res, err := petRepo.GetPetProfileByID(1000, 10000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestUpdatePetProfile(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	newPet := entity.Pet{}
	newPet.Name = "pet1"
	newPet.UserID = 1

	db.Create(&newPet)

	petRepo := NewPetRepository(db)

	t.Run("1. Success update pet", func(t *testing.T) {
		updatedPet := entity.Pet{}
		updatedPet.Name = "pet1 new"
		updatedPet.UserID = 1

		res, err := petRepo.UpdatePetProfile(1, 1, updatedPet)

		assert.Nil(t, err)
		assert.Equal(t, "pet1 new", res.Name)
	})

	t.Run("2. Error update pet - pet not found", func(t *testing.T) {
		updatedPet := entity.Pet{}
		updatedPet.Name = "pet1 new"
		updatedPet.UserID = 1

		res, err := petRepo.UpdatePetProfile(1000, 1000, updatedPet)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestDeletePet(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	newPet := entity.Pet{}
	newPet.Name = "pet1"
	newPet.UserID = 1

	db.Create(&newPet)

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
