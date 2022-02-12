package category

import (
	"errors"
	"petshop/entity"

	"gorm.io/gorm"
)

type Category interface {
	CreateCategory(newCategory entity.Category) (entity.Category, error)
	GetAllCategory() ([]entity.Category, error)
	GetCategoryByID(categoryID int) (entity.Category, error)
	UpdateCategory(categoryID int, updatedCategory entity.Category) (entity.Category, error)
	DeleteCategory(categoryID int) (entity.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (cr *categoryRepository) CreateCategory(newCategory entity.Category) (entity.Category, error) {
	err := cr.db.Save(&newCategory).Error

	if err != nil {
		return newCategory, err
	}

	return newCategory, nil
}

func (cr *categoryRepository) GetAllCategory() ([]entity.Category, error) {
	categories := []entity.Category{}

	err := cr.db.Find(&categories).Error

	if err != nil || len(categories) == 0 {
		return categories, errors.New("category not found")
	}

	return categories, err
}

func (cr *categoryRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	category := entity.Category{}

	err := cr.db.Where("id = ?", categoryID).Find(&category).Error

	if err != nil || category.ID == 0 {
		return category, errors.New("category not found")
	}

	return category, err
}

func (cr *categoryRepository) UpdateCategory(categoryID int, updatedCategory entity.Category) (entity.Category, error) {
	category := entity.Category{}

	err := cr.db.Where("id = ?", categoryID).Find(&category).Error

	if err != nil || category.ID == 0 {
		return category, errors.New("category not found")
	}

	cr.db.Model(&category).Updates(&updatedCategory)

	return updatedCategory, err
}

func (cr *categoryRepository) DeleteCategory(categoryID int) (entity.Category, error) {
	category := entity.Category{}

	err := cr.db.Where("id = ?", categoryID).Find(&category).Error

	if err != nil || category.ID == 0 {
		return category, errors.New("category not found")
	}

	cr.db.Delete(&category)

	return category, err
}
