package store

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
)

type Store interface {
	FindCityByID(cityID int) (entity.City, error)
	CreateStore(newStore entity.Store) (entity.Store, error)
	GetAllStoreByUserID(userID int) ([]entity.Store, error)
	GetStoreProfile(storeID, userID int) (entity.Store, error)
	UpdateStoreProfile(storeID, userID int, updatedStore entity.Store) (entity.Store, error)
	DeleteStore(storeID, userID int) (entity.Store, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *storeRepository {
	return &storeRepository{db}
}

func (sr *storeRepository) FindCityByID(cityID int) (entity.City, error) {
	city := entity.City{}

	err := sr.db.Where("id = ?", cityID).Find(&city).Error

	if err != nil || city.ID == 0 {
		return city, err
	}

	return city, nil
}

func (sr *storeRepository) CreateStore(newStore entity.Store) (entity.Store, error) {
	err := sr.db.Save(&newStore).Error

	if err != nil {
		return newStore, err
	}

	return newStore, nil
}

func (sr *storeRepository) GetAllStoreByUserID(userID int) ([]entity.Store, error) {
	stores := []entity.Store{}

	err := sr.db.Where("user_id = ?", userID).Find(&stores).Error

	if err != nil || len(stores) == 0 {
		return stores, errors.New("store not found")
	}

	return stores, nil
}

func (sr *storeRepository) GetStoreProfile(storeID, userID int) (entity.Store, error) {
	store := entity.Store{}

	err := sr.db.Where("id = ? AND user_id = ?", storeID, userID).Find(&store).Error

	if err != nil || store.ID == 0 {
		return store, errors.New("store not found")
	}

	return store, nil
}

func (sr *storeRepository) UpdateStoreProfile(storeID, userID int, updatedStore entity.Store) (entity.Store, error) {
	store := entity.Store{}

	err := sr.db.Where("id = ? AND user_id = ?", storeID, userID).Find(&store).Error

	if err != nil || store.ID == 0 {
		return store, errors.New("store not found")
	}

	sr.db.Model(&store).Updates(&updatedStore)

	return store, nil
}

func (sr *storeRepository) DeleteStore(storeID, userID int) (entity.Store, error) {
	store := entity.Store{}

	err := sr.db.Where("id = ? AND user_id = ?", storeID, userID).Find(&store).Error

	if err != nil || store.ID == 0 {
		return store, errors.New("store not found")
	}

	sr.db.Delete(&store)

	return store, nil
}
