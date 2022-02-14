package category

import (
	"gorm.io/gorm"
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.Migrator().DropTable(&entity.Category{})
	db.AutoMigrate(&entity.Category{})
}

func TestCreateCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	categoryRepo := NewCategoryRepository(db)

	newCategory := entity.Category{}
	newCategory.Name = "category2"

	t.Run(
		"1. Success create category", func(t *testing.T) {
			res, err := categoryRepo.CreateCategory(newCategory)

			assert.Nil(t, err)
			assert.Equal(t, "category2", res.Name)
		},
	)

	t.Run(
		"2. Error create category - Duplicate name", func(t *testing.T) {
			res, err := categoryRepo.CreateCategory(newCategory)

			assert.NotNil(t, err)
			assert.Equal(t, 0, int(res.ID))
		},
	)
}

func TestGetAllCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	categoryRepo := NewCategoryRepository(db)

	t.Run(
		"1. Success get category", func(t *testing.T) {
			res, err := categoryRepo.GetAllCategory()

			assert.Nil(t, err)
			assert.Equal(t, "category2", res[0].Name)
		},
	)

	t.Run(
		"2. Error get category", func(t *testing.T) {
			category1 := entity.Category{}
			db.Where("id = ?", 1).Delete(&category1)
			res, err := categoryRepo.GetAllCategory()

			assert.NotNil(t, err)
			assert.Equal(t, false, len(res) > 0)
		},
	)
}

func TestGetCategoryByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	categoryRepo := NewCategoryRepository(db)
	newCategory := entity.Category{
		Name: "Makanan",
	}

	db.Create(&newCategory)

	t.Run(
		"1. Success get category", func(t *testing.T) {

			res, err := categoryRepo.GetCategoryByID(3)

			assert.Nil(t, err)
			assert.Equal(t, "Makanan", res.Name)
		},
	)

	t.Run(
		"2. Error get category - category id not found", func(t *testing.T) {
			res, err := categoryRepo.GetCategoryByID(1000)

			assert.NotNil(t, err)
			assert.Equal(t, "", res.Name)
		},
	)
}

func TestUpdateCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	categoryRepo := NewCategoryRepository(db)

	t.Run(
		"1. Success update category", func(t *testing.T) {
			mockRequest := entity.Category{
				Model: gorm.Model{},
				Name:  "Updated Name",
			}
			res, err := categoryRepo.UpdateCategory(3, mockRequest)

			assert.Nil(t, err)
			assert.Equal(t, "Updated Name", res.Name)
		},
	)

	t.Run(
		"2. Error update category - category not found", func(t *testing.T) {
			mockRequest := entity.Category{
				Model: gorm.Model{},
				ID:    100,
				Name:  "Updated Name",
			}
			res, err := categoryRepo.UpdateCategory(1000, mockRequest)

			assert.NotNil(t, err)
			assert.Equal(t, "", res.Name)
		},
	)
}

func TestDeleteCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	categoryRepo := NewCategoryRepository(db)

	t.Run(
		"1. Success Delete category", func(t *testing.T) {
			_, err := categoryRepo.DeleteCategory(3)

			assert.Nil(t, err)
		},
	)

	t.Run(
		"2. Error update category - category not found", func(t *testing.T) {
			_, err := categoryRepo.DeleteCategory(100)

			assert.NotNil(t, err)
		},
	)
}
