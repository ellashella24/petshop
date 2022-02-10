package city

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
)

type City interface {
	CreateCity(newCity entity.City) (entity.City, error)
	GetAllCity() ([]entity.City, error)
	GetCityByID(cityID int) (entity.City, error)
	UpdateCity(cityID int, updatedCity entity.City) (entity.City, error)
	DeleteCity(cityID int) (entity.City, error)
}

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *cityRepository {
	return &cityRepository{db}
}

func (cr *cityRepository) CreateCity(newCity entity.City) (entity.City, error) {
	err := cr.db.Save(&newCity).Error

	if err != nil {
		return newCity, err
	}

	return newCity, nil
}

func (cr *cityRepository) GetAllCity() ([]entity.City, error) {
	cities := []entity.City{}

	err := cr.db.Find(&cities).Error

	if err != nil || len(cities) == 0 {
		return cities, errors.New("city not found")
	}

	return cities, err
}

func (cr *cityRepository) GetCityByID(cityID int) (entity.City, error) {
	city := entity.City{}

	err := cr.db.Where("id = ?", cityID).Find(&city).Error

	if err != nil || city.ID == 0 {
		return city, errors.New("city not found")
	}

	return city, err
}

func (cr *cityRepository) UpdateCity(cityID int, updatedCity entity.City) (entity.City, error) {
	city := entity.City{}

	err := cr.db.Where("id = ?", cityID).Find(&city).Error

	if err != nil || city.ID == 0 {
		return city, errors.New("city not found")
	}

	cr.db.Model(&city).Updates(&updatedCity)

	return updatedCity, err
}

func (cr *cityRepository) DeleteCity(cityID int) (entity.City, error) {
	city := entity.City{}

	err := cr.db.Where("id = ?", cityID).Find(&city).Error

	if err != nil || city.ID == 0 {
		return city, errors.New("city not found")
	}

	cr.db.Delete(&city)

	return city, err
}
