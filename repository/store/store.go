package store

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
)

type Store interface {
	FindCityByID(cityID int) (entity.City, error)
	CreateStore(newStore entity.Store) (entity.Store, error)
	GetAllStore() ([]entity.Store, error)
	GetStoreProfile(storeID, userID int) (entity.Store, error)
	UpdateStoreProfile(storeID, userID int, updatedStore entity.Store) (entity.Store, error)
	DeleteStore(storeID, userID int) (entity.Store, error)
	GetListTransactionByStoreID(storeID int) ([]entity.Transaction, []entity.TransactionDetail, []entity.Product, error)
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

func (sr *storeRepository) GetAllStore() ([]entity.Store, error) {
	stores := []entity.Store{}

	err := sr.db.Find(&stores).Error

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

func (sr *storeRepository) GetListTransactionByStoreID(storeID int) ([]entity.Transaction, []entity.TransactionDetail, []entity.Product, error) {
	var transactionDetail []entity.TransactionDetail
	var transactions []entity.Transaction
	var products []entity.Product

	err := sr.db.Table("transaction_details").Joins("join products on transaction_details.product_id = products.id").Where("products.store_id = ?", storeID).Find(&transactionDetail).Error

	if err != nil || len(transactionDetail) == 0 {
		return transactions, transactionDetail, products, errors.New("not found")
	}

	for i := range transactionDetail {
		var transaction entity.Transaction
		var product entity.Product

		sr.db.Where("id", transactionDetail[i].TransactionID).First(&transaction)

		sr.db.Where("id", transactionDetail[i].ProductID).First(&product)

		if len(transactions) == 0 {
			transactions = append(transactions, transaction)
		}

		if transactions[len(transactions)-1].ID != transaction.ID {
			transactions = append(transactions, transaction)
		}

		products = append(products, product)
	}

	return transactions, transactionDetail, products, nil
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
