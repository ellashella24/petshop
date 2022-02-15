package user

import (
	"fmt"
	"petshop/config"
	"petshop/entity"
	"petshop/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSetup(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)
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

	city := entity.City{}
	city.Name = "city1"

	db.Create(&city)

}

func TestCreateUser(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	userRepo := NewUserRepository(db)

	t.Run(
		"1. Success create user", func(t *testing.T) {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)

			mockUser := entity.User{}
			mockUser.Name = "user1"
			mockUser.Email = "user1@mail.com"
			mockUser.Password = string(hashedPassword)
			mockUser.CityID = 1

			res, err := userRepo.CreateUser(mockUser)

			assert.Nil(t, err)
			assert.Equal(t, mockUser.Name, res.Name)
		},
	)

	t.Run(
		"2. Error create user - Duplicate email", func(t *testing.T) {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)

			mockUser := entity.User{}
			mockUser.Name = "user1"
			mockUser.Email = "user1@mail.com"
			mockUser.Password = string(hashedPassword)
			mockUser.CityID = 1

			res, err := userRepo.CreateUser(mockUser)

			assert.NotNil(t, err)
			assert.Equal(t, 0, int(res.ID))
		},
	)
}

func TestGetUserByEmail(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)
	user := entity.User{}
	user.Name = "user1"
	user.Email = "user1@mail.com"
	user.Password = string(hashedPassword)
	user.CityID = 1

	db.Create(&user)

	userRepo := NewUserRepository(db)

	t.Run(
		"1. Success get user", func(t *testing.T) {
			res, err := userRepo.GetUserByEmail("user1@mail.com")

			assert.Nil(t, err)
			assert.Equal(t, "user1@mail.com", res.Email)
		},
	)

	t.Run(
		"2. Error get user - email not found", func(t *testing.T) {
			res, err := userRepo.GetUserByEmail("user2222@mail.com")

			fmt.Println(res)

			assert.NotNil(t, err)
			assert.Equal(t, "", res.Email)
		},
	)
}

func TestGetUserByID(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.MinCost)
	newUser := entity.User{}
	newUser.Name = "user1"
	newUser.Email = "user1@mail.com"
	newUser.Password = string(hashedPassword)
	newUser.CityID = 1

	db.Save(&newUser)

	userRepo := NewUserRepository(db)

	t.Run(
		"1. Success get user", func(t *testing.T) {
			res, err := userRepo.GetUserByID(1)

			assert.Nil(t, err)
			assert.Equal(t, 1, int(res.ID))
		},
	)

	t.Run(
		"2. Error get user - id not found", func(t *testing.T) {
			res, _ := userRepo.GetUserByID(100)

			assert.Equal(t, 0, int(res.ID))
		},
	)
}

func TestUpdateUser(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	userRepo := NewUserRepository(db)

	t.Run(
		"1. Success update user", func(t *testing.T) {

			mockRequest := entity.User{
				Name: "User 1 Update",
			}
			res, err := userRepo.UpdateUser(1, mockRequest)

			assert.Nil(t, err)
			assert.Equal(t, "User 1 Update", res.Name)

		},
	)

	t.Run(
		"2. Error update user - id not found", func(t *testing.T) {
			mockRequest := entity.User{
				Name: "User 1 Update",
			}
			res, _ := userRepo.UpdateUser(100, mockRequest)

			assert.Equal(t, "", res.Name)
		},
	)
}

func TestDeleteUser(t *testing.T) {
	config := config.GetConfig()
	db := util.InitDB(config)

	userRepo := NewUserRepository(db)

	t.Run(
		"1. Success delete user", func(t *testing.T) {
			_, err := userRepo.DeleteUser(1)

			assert.Nil(t, err)
		},
	)

	t.Run(
		"2. Error delete user", func(t *testing.T) {

			_, err := userRepo.DeleteUser(100)

			assert.NotNil(t, err)
		},
	)
}
