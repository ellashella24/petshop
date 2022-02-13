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
	GetListTransactionByStoreID(storeID int) ([]entity.Transaction, error)
	GetGroomingStatusByPetID(petID, storeID int) (entity.GroomingStatus, error)
	UpdateGroomingStatus(petID, storeID int) (entity.GroomingStatus, error)
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

func (sr *storeRepository) GetListTransactionByStoreID(storeID int) ([]entity.Transaction, error) {
	transactions := []entity.Transaction{}

	sr.db.Select("transactions.invoice_id, transactions.user_id, transactions.total_price, transactions.payment_status, transactions.paid_at, transactions.payment_method").Table("transactions").Joins("join transaction_details ON transactions.id = transaction_details.transaction_id").Joins("join products ON transaction_details.product_id = products.id").Joins("join stores ON products.store_id = stores.id").Where("stores.id = ?", storeID).Group("transactions.invoice_id").Find(&transactions)

	if len(transactions) == 0 {
		return transactions, errors.New("transaction not found")
	}

	return transactions, nil
}

func (sr *storeRepository) UpdateGroomingStatus(petID, storeID int) (entity.GroomingStatus, error) {
	grooming_status := entity.GroomingStatus{}

	sr.db.Joins("join transaction_details ON grooming_statuses.transaction_detail_id = transaction_details.id").Joins("join products ON transaction_details.product_id = products.id").Where("grooming_statuses.pet_id = ? AND products.store_id = ? AND products.category_id = ?", petID, storeID, 1).Find(&grooming_status)

	if grooming_status.Status == "" {
		return grooming_status, errors.New("not found grooming status")
	} else if grooming_status.Status == "TELAH DIBAYAR" {
		sr.db.Model(&grooming_status).Update("status", "PROSES PICKUP")
	} else if grooming_status.Status == "PROSES PICKUP" {
		sr.db.Model(&grooming_status).Update("status", "PROSES GROOMING")
	} else if grooming_status.Status == "PROSES GROOMING" {
		sr.db.Model(&grooming_status).Update("status", "DELIVERY KE USER")
	} else {
		return grooming_status, errors.New("can't update grooming status")
	}

	return grooming_status, nil
}

func (sr *storeRepository) GetGroomingStatusByPetID(petID, storeID int) (entity.GroomingStatus, error) {
	grooming_status := entity.GroomingStatus{}

	sr.db.Joins("join transaction_details ON grooming_statuses.transaction_detail_id = transaction_details.id").Joins("join products ON transaction_details.product_id = products.id").Where("grooming_statuses.pet_id = ? AND products.store_id = ? AND products.category_id = ?", petID, storeID, 1).Find(&grooming_status)

	if grooming_status.Status == "" {
		return grooming_status, errors.New("not found grooming status")
	}

	return grooming_status, nil
}
