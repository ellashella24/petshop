package transaction

import (
	"gorm.io/gorm"
	"petshop/entity"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

type Transaction interface {
	Transaction(newTransactions entity.Transaction) (entity.Transaction, error)
	TransactionDetail(newDetailTransactions entity.TransactionDetail) (entity.Transaction, error)
}

func (tr *TransactionRepository) Transaction(newTransactions entity.Transaction) (entity.Transaction, error) {
	err := tr.db.Save(&newTransactions).Error
	if err != nil {
		return newTransactions, err
	}
	return newTransactions, nil
}

func (tr *TransactionRepository) TransactionDetail(newDetailTransactions entity.TransactionDetail) (
	entity.TransactionDetail, error,
) {
	err := tr.db.Save(&newDetailTransactions).Error
	if err != nil {
		return newDetailTransactions, err
	}
	return newDetailTransactions, nil
}
