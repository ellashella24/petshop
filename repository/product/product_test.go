package product

import (
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db = util.InitDB(config)

	db.Migrator().DropTable(&entity.TransactionDetail{})
	db.Migrator().DropTable(&entity.Transaction{})
	db.Migrator().DropTable(&entity.GroomingStatus{})
	db.Migrator().DropTable(&entity.StockHistory{})
	db.Migrator().DropTable(&entity.Cart{})
	db.Migrator().DropTable(&entity.Product{})
	db.Migrator().DropTable(&entity.Category{})
	db.Migrator().DropTable(&entity.Store{})
	db.Migrator().DropTable(&entity.User{})
	db.Migrator().DropTable(&entity.City{})
	db.Migrator().DropTable(&entity.Pet{})

	db.AutoMigrate(&entity.City{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Pet{})
	db.AutoMigrate(&entity.Store{})
	db.AutoMigrate(&entity.Category{})
	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.Cart{})
	db.AutoMigrate(&entity.StockHistory{})
	db.AutoMigrate(&entity.GroomingStatus{})
	db.AutoMigrate(&entity.Transaction{})
	db.AutoMigrate(&entity.TransactionDetail{})

	newCity := entity.City{
		Name: "City 1",
	}

	db.Create(&newCity)

	newCategory := entity.Category{
		Name: "Grooming",
	}

	db.Create(&newCategory)

	newCategory1 := entity.Category{
		Name: "Makanan",
	}

	db.Create(&newCategory1)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)

	newUser := entity.User{
		Name:     "User 1",
		Email:    "user1@mail.com",
		Password: string(hashedPassword),
		CityID:   1,
	}

	db.Create(&newUser)

	newStore := entity.Store{
		Name:   "Store 1",
		UserID: 1,
		CityID: 1,
	}

	db.Create(&newStore)
}

func TestCreateProduct(t *testing.T) {
	productRepo := NewProductRepository(db)

	t.Run(
		"1. Success Create Product", func(t *testing.T) {
			newProduct := entity.Product{
				Name:       "Whiskas",
				Price:      100000,
				Stock:      100,
				StoreID:    1,
				CategoryID: 2,
			}

			res, err := productRepo.CreateProduct(1, newProduct)

			assert.Nil(t, err)
			assert.Equal(t, "Whiskas", res.Name)
		},
	)

	t.Run(
		"2. Error Create Product", func(t *testing.T) {
			newProduct := entity.Product{
				Name:       "Whiskas",
				Price:      100000,
				Stock:      100,
				StoreID:    0,
				CategoryID: 100000,
			}

			_, err := productRepo.CreateProduct(10, newProduct)

			assert.NotNil(t, err)
		},
	)
}

func TestGetAllProduct(t *testing.T) {
	productRepo := NewProductRepository(db)

	t.Run(
		"1. Success get all product", func(t *testing.T) {
			res, err := productRepo.GetAllProduct()

			assert.Nil(t, err)
			assert.Equal(t, "Whiskas", res[0].Name)
		},
	)

	t.Run(
		"2. Error get all product", func(t *testing.T) {
			productRepo.DeleteProduct(1)

			res, err := productRepo.GetAllProduct()

			assert.NotNil(t, err)
			assert.Equal(t, []entity.Product{}, res)
		},
	)
}

func TestGetProductByID(t *testing.T) {
	productRepo := NewProductRepository(db)

	newProduct := entity.Product{
		Name:       "Excel",
		Price:      100000,
		Stock:      100,
		StoreID:    1,
		CategoryID: 2,
	}

	db.Create(&newProduct)

	// fmt.Println(productRepo.GetAllProduct())

	t.Run(
		"1. Success get product", func(t *testing.T) {
			res, err := productRepo.GetProductByID(2)

			assert.Nil(t, err)
			assert.Equal(t, "Excel", res.Name)
		},
	)

	t.Run(
		"2. Error get product", func(t *testing.T) {
			res, err := productRepo.GetProductByID(10000)

			assert.NotNil(t, err)
			assert.Equal(t, "", res.Name)
		},
	)
}

func TestGetProductByStoreID(t *testing.T) {
	productRepo := NewProductRepository(db)

	t.Run(
		"1. Success get product", func(t *testing.T) {
			res, err := productRepo.GetProductByStoreID(1)

			assert.Nil(t, err)
			assert.Equal(t, "Excel", res[0].Name)
		},
	)

	t.Run(
		"2. Error get product", func(t *testing.T) {
			res, err := productRepo.GetProductByStoreID(10000)

			assert.NotNil(t, err)
			assert.Equal(t, []entity.Product([]entity.Product{}), res)
		},
	)
}

func TestGetStockHistory(t *testing.T) {
	productRepo := NewProductRepository(db)

	t.Run(
		"1. Success get product stock history ", func(t *testing.T) {
			res, err := productRepo.GetStockHistory(1)

			assert.Nil(t, err)
			assert.Equal(t, 100, res[0].Stock)
			assert.Equal(t, 1, int(res[0].ProductID))
		},
	)

	t.Run(
		"2. Error get product stock hsitory", func(t *testing.T) {
			res, err := productRepo.GetStockHistory(10000)

			assert.NotNil(t, err)
			assert.Equal(t, []entity.StockHistory([]entity.StockHistory{}), res)
		},
	)
}

func TestUpdateProduct(t *testing.T) {
	productRepo := NewProductRepository(db)

	t.Run(
		"1. Success update product", func(t *testing.T) {
			updatedProduct := entity.Product{
				Name:       "Me-o",
				Price:      100000,
				Stock:      100,
				StoreID:    1,
				CategoryID: 2,
			}

			res, err := productRepo.UpdateProduct(2, updatedProduct)

			assert.Nil(t, err)
			assert.Equal(t, "Me-o", res.Name)
		},
	)

	t.Run(
		"2. Error update product", func(t *testing.T) {
			updatedProduct := entity.Product{
				Name:       "Me-o",
				Price:      100000,
				Stock:      100,
				StoreID:    1,
				CategoryID: 2000000,
			}

			res, err := productRepo.UpdateProduct(10000, updatedProduct)

			assert.NotNil(t, err)
			assert.Equal(t, "", res.Name)
		},
	)
}

func TestDeleteProduct(t *testing.T) {
	productRepo := NewProductRepository(db)

	t.Run(
		"1. Success get product", func(t *testing.T) {
			_, err := productRepo.DeleteProduct(2)

			assert.Nil(t, err)
		},
	)

	t.Run(
		"2. Error get product", func(t *testing.T) {
			_, err := productRepo.DeleteProduct(10000)

			assert.NotNil(t, err)
		},
	)
}
