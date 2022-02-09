package helper

import (
	"errors"
	"petshop/entity"
)

func GroomingStatusHelper(petID int, transactionDetailID uint) error {
	var groomingStatus entity.GroomingStatus
	groomingStatus = entity.GroomingStatus{
		PetID:               uint(petID),
		TransactionDetailID: transactionDetailID,
	}
	err := db.Save(&groomingStatus).Error

	if err != nil {
		return errors.New("error")
	}

	return nil
}
