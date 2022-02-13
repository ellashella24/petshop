package pet

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
)

type Pet interface {
	CreatePet(newPet entity.Pet) (entity.Pet, error)
	GetAllPetByUserID(userID int) ([]entity.Pet, error)
	GetPetProfileByID(petID, userID int) (entity.Pet, error)
	UpdatePetProfile(petID, userID int, updatedPet entity.Pet) (entity.Pet, error)
	DeletePet(petID, userID int) (entity.Pet, error)
	GetGroomingStatusByPetID(petID, userID int) (entity.GroomingStatus, error)
	UpdateFinalGroomingStatus(petID, userID int) (entity.GroomingStatus, error)
}

type petRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) *petRepository {
	return &petRepository{db}
}

func (pr *petRepository) CreatePet(newPet entity.Pet) (entity.Pet, error) {
	err := pr.db.Save(&newPet).Error

	if err != nil {
		return newPet, err
	}

	return newPet, nil
}

func (pr *petRepository) GetAllPetByUserID(userID int) ([]entity.Pet, error) {
	pets := []entity.Pet{}

	err := pr.db.Where("user_id = ?", userID).Find(&pets).Error

	if err != nil || len(pets) == 0 {
		return pets, errors.New("pet not found")
	}

	return pets, err
}

func (pr *petRepository) GetPetProfileByID(petID, userID int) (entity.Pet, error) {
	pet := entity.Pet{}

	err := pr.db.Where("id = ? AND user_id = ?", petID, userID).Find(&pet).Error

	if err != nil || pet.ID == 0 {
		return pet, errors.New("pet not found")
	}

	return pet, err
}

func (pr *petRepository) UpdatePetProfile(petID, userID int, updatedPet entity.Pet) (entity.Pet, error) {
	pet := entity.Pet{}

	err := pr.db.Where("id = ? AND user_id = ?", petID, userID).Find(&pet).Error

	if err != nil || pet.ID == 0 {
		return pet, errors.New("pet not found")
	}

	pr.db.Model(&pet).Updates(&updatedPet)

	return updatedPet, err
}

func (pr *petRepository) DeletePet(petID, userID int) (entity.Pet, error) {
	pet := entity.Pet{}

	err := pr.db.Where("id = ? AND user_id = ?", petID, userID).Find(&pet).Error

	if err != nil || pet.ID == 0 {
		return pet, errors.New("pet not found")
	}

	pr.db.Delete(&pet)

	return pet, err
}

func (pr *petRepository) GetGroomingStatusByPetID(petID, userID int) (entity.GroomingStatus, error) {
	grooming_status := entity.GroomingStatus{}

	pr.db.Joins("join pets ON grooming_statuses.pet_id = pets.id").Where("pets.id = ? AND pets.user_id = ?", petID, userID).Find(&grooming_status)

	if grooming_status.Status == "" {
		return grooming_status, errors.New("not found grooming status")
	}

	return grooming_status, nil
}

func (pr *petRepository) UpdateFinalGroomingStatus(petID, userID int) (entity.GroomingStatus, error) {
	grooming_status := entity.GroomingStatus{}

	pr.db.Joins("join pets ON grooming_statuses.pet_id = pets.id").Where("pets.id = ? AND pets.user_id = ?", petID, userID).Find(&grooming_status)

	if grooming_status.Status == "" {
		return grooming_status, errors.New("not found grooming status")
	} else if grooming_status.Status == "DELIVERY KE USER" {
		pr.db.Model(&grooming_status).Update("status", "SELESAI")
	} else {
		return grooming_status, errors.New("can't update grooming status")
	}

	return grooming_status, nil
}
