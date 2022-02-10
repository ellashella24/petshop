package category

import (
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

	t.Run("1. Success create category", func(t *testing.T) {
		res, err := categoryRepo.CreateCategory(newCategory)

		assert.Nil(t, err)
		assert.Equal(t, "category2", res.Name)
	})

	t.Run("2. Error create category - Duplicate name", func(t *testing.T) {
		res, err := categoryRepo.CreateCategory(newCategory)

		assert.NotNil(t, err)
		assert.Equal(t, 0, int(res.ID))
	})
}

func TestGetAllCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	categoryRepo := NewCategoryRepository(db)

	newCategory := entity.Category{}
	newCategory.Name = "category2"

	db.Create(&newCategory)

	t.Run("1. Success get category", func(t *testing.T) {
		res, err := categoryRepo.GetAllCategory()

		assert.Nil(t, err)
		assert.Equal(t, "category2", res[1].Name)
	})

	t.Run("2. Error get category", func(t *testing.T) {
		category1 := entity.Category{}
		db.Where("id = ?", 1).Delete(&category1)
		category2 := entity.Category{}
		db.Where("id = ?", 2).Delete(&category2)
		res, err := categoryRepo.GetAllCategory()

		assert.NotNil(t, err)
		assert.Equal(t, false, len(res) > 0)
	})
}

func TestGetCategoryByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.AutoMigrate(&entity.Category{})

	categoryRepo := NewCategoryRepository(db)

	newCategory := entity.Category{}
	newCategory.Name = "category2"

	db.Create(&newCategory)

	t.Run("1. Success get category", func(t *testing.T) {
		res, err := categoryRepo.GetCategoryByID(2)

		assert.Nil(t, err)
		assert.Equal(t, newCategory.Name, res.Name)
	})

	t.Run("2. Error get category - category id not found", func(t *testing.T) {
		res, err := categoryRepo.GetCategoryByID(1000)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestUpdateCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.AutoMigrate(&entity.Category{})

	categoryRepo := NewCategoryRepository(db)

	newCategory := entity.Category{}
	newCategory.Name = "category2"

	db.Create(&newCategory)

	updatedCategory := entity.Category{}
	updatedCategory.Name = "category2 new"

	t.Run("1. Success update category", func(t *testing.T) {
		res, err := categoryRepo.UpdateCategory(2, updatedCategory)

		assert.Nil(t, err)
		assert.Equal(t, updatedCategory.Name, res.Name)
	})

	t.Run("2. Error update category - category not found", func(t *testing.T) {
		res, err := categoryRepo.UpdateCategory(1000, updatedCategory)

		assert.NotNil(t, err)
		assert.Equal(t, "", res.Name)
	})
}

func TestDeleteCategory(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	db.AutoMigrate(&entity.Category{})

	categoryRepo := NewCategoryRepository(db)

	newCategory := entity.Category{}
	newCategory.Name = "category2"

	db.Create(&newCategory)

	t.Run("1. Success update category", func(t *testing.T) {
		_, err := categoryRepo.DeleteCategory(2)

		assert.Nil(t, err)
	})

	t.Run("2. Error update category - category not found", func(t *testing.T) {
		_, err := categoryRepo.DeleteCategory(100)

		assert.NotNil(t, err)
	})
}
