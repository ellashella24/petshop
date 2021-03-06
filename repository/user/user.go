package user

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
)

type User interface {
	CreateUser(newUser entity.User) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUserByID(userID int) (entity.User, error)
	UpdateUser(userID int, updatedUser entity.User) (entity.User, error)
	DeleteUser(userID int) (entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(newUser entity.User) (entity.User, error) {
	err := ur.db.Save(&newUser).Error

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (ur *userRepository) GetUserByEmail(email string) (entity.User, error) {
	user := entity.User{}

	err := ur.db.Where("email = ?", email).Find(&user).Error

	if err != nil || user.ID == 0 {
		return user, errors.New("email not found")
	}

	return user, nil
}

func (ur *userRepository) GetUserByID(userID int) (entity.User, error) {
	user := entity.User{}

	err := ur.db.Where("id = ?", userID).Find(&user).Error

	if err != nil || user.ID == 0 {
		return user, errors.New("user not found")
	}

	return user, nil
}

func (ur *userRepository) UpdateUser(userID int, updatedUser entity.User) (entity.User, error) {
	user := entity.User{}

	err := ur.db.Where("id = ?", userID).Find(&user).Error

	if err != nil || user.ID == 0 {
		return user, errors.New("user not found")
	}

	ur.db.Model(&user).Updates(updatedUser)

	return updatedUser, nil
}

func (ur *userRepository) DeleteUser(userID int) (entity.User, error) {
	user := entity.User{}

	err := ur.db.Where("id = ?", userID).Find(&user).Error

	if err != nil || user.ID == 0 {
		return user, errors.New("user not found")
	}

	ur.db.Delete(&user)

	return user, nil
}
