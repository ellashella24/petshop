package city

import (
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.Migrator().DropTable(&entity.City{})
	db.AutoMigrate(&entity.City{})
}

func TestCreateCity(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cityRepo := NewCityRepository(db)

	newCity := entity.City{}
	newCity.Name = "city2"

	t.Run("1. Success create city", func(t *testing.T) {
		res, err := cityRepo.CreateCity(newCity)

		assert.Nil(t, err)
		assert.Equal(t, "city2", res.Name)
	})

	t.Run("2. Error create city - Duplicate name", func(t *testing.T) {
		res, err := cityRepo.CreateCity(newCity)

		assert.NotNil(t, err)
		assert.Equal(t, 0, int(res.ID))
	})
}

func TestGetAllCity(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	cityRepo := NewCityRepository(db)

	newCity := entity.City{}
	newCity.Name = "city2"

	db.Create(&newCity)

	t.Run("1. Success get city", func(t *testing.T) {
		res, err := cityRepo.GetAllCity()

		assert.Nil(t, err)
		assert.Equal(t, "city2", res[1].Name)
	})

	t.Run("2. Error get city", func(t *testing.T) {
		city1 := entity.City{}
		db.Where("id = ?", 1).Delete(&city1)
		city2 := entity.City{}
		db.Where("id = ?", 2).Delete(&city2)
		res, err := cityRepo.GetAllCity()

		assert.NotNil(t, err)
		assert.Equal(t, false, len(res) > 0)
	})
}

func TestGetCityByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.AutoMigrate(&entity.City{})

	cityRepo := NewCityRepository(db)

	newCity := entity.City{}
	newCity.Name = "city2"

	db.Create(&newCity)

	t.Run("1. Success get city", func(t *testing.T) {
		res, err := cityRepo.GetCityByID(2)

		assert.Nil(t, err)
		assert.Equal(t, newCity.Name, res.Name)
	})

	t.Run("2. Error get city - city id not found", func(t *testing.T) {
		res, err := cityRepo.GetCityByID(1000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestUpdateCity(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.AutoMigrate(&entity.City{})

	cityRepo := NewCityRepository(db)

	newCity := entity.City{}
	newCity.Name = "city2"

	db.Create(&newCity)

	updatedCity := entity.City{}
	updatedCity.Name = "city2 new"

	t.Run("1. Success update city", func(t *testing.T) {
		res, err := cityRepo.UpdateCity(2, updatedCity)

		assert.Nil(t, err)
		assert.Equal(t, updatedCity.Name, res.Name)
	})

	t.Run("2. Error update city - city not found", func(t *testing.T) {
		res, err := cityRepo.UpdateCity(1000, updatedCity)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestDeleteCity(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.AutoMigrate(&entity.City{})

	cityRepo := NewCityRepository(db)

	newCity := entity.City{}
	newCity.Name = "city2"

	db.Create(&newCity)

	t.Run("1. Success update city", func(t *testing.T) {
		_, err := cityRepo.DeleteCity(2)

		assert.Nil(t, err)
	})

	t.Run("2. Error update city - city not found", func(t *testing.T) {
		_, err := cityRepo.DeleteCity(100)

		assert.NotNil(t, err)
	})
}
