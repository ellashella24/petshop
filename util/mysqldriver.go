package util

import (
	"fmt"
	"petshop/config"
	"petshop/entity"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(config *config.AppConfig) *gorm.DB {
	connectionString := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Address,
		config.DB_Port,
		config.Name,
	)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Info("failed to connect database :", err)
		panic(err)
	}

	InitMigrate(db)

	return db
}

func InitMigrate(db *gorm.DB) {
	db.Migrator().DropTable(&entity.TransactionDetail{})
	db.Migrator().DropTable(&entity.Transaction{})
	db.Migrator().DropTable(&entity.GroomingStatus{})
	db.Migrator().DropTable(&entity.StockHistory{})
	db.Migrator().DropTable(&entity.Cart{})
	db.Migrator().DropTable(&entity.Product{})
	db.Migrator().DropTable(&entity.Category{})
	db.Migrator().DropTable(&entity.Store{})
	db.Migrator().DropTable(&entity.Pet{})
	db.Migrator().DropTable(&entity.User{})
	db.Migrator().DropTable(&entity.City{})

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
	city.Name = "admin city"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)
	admin := entity.User{}
	admin.Name = "admin1"
	admin.Email = "admin@admin.com"
	admin.Password = string(hashedPassword)
	admin.CityID = 1
	admin.Role = "admin"

	user := entity.User{}
	user.Name = "user"
	user.Email = "user@kruay.com"
	user.Password = string(hashedPassword)
	user.CityID = 1
	user.Role = "user"

	store := entity.Store{}
	store.Name = "store1"
	store.UserID = 1
	store.CityID = 1

	store1 := entity.Store{}
	store1.Name = "store2"
	store1.UserID = 2
	store1.CityID = 1

	category := entity.Category{
		Name:    "Grooming",
		Product: nil,
	}
	category1 := entity.Category{
		Name:    "Makanan",
		Product: nil,
	}
	pet := entity.Pet{
		Name:   "Kucing",
		UserID: 2,
	}
	pet2 := entity.Pet{
		Name:   "Kucing",
		UserID: 2,
	}
	product := entity.Product{
		Name:       "Whiskas",
		Price:      10000,
		Stock:      10,
		StoreID:    1,
		CategoryID: 2,
		Category:   entity.Category{},
	}
	product1 := entity.Product{
		Name:       "Shampoan",
		Price:      10000,
		StoreID:    2,
		CategoryID: 1,
		Category:   entity.Category{},
	}

	db.Create(&city)
	db.Create(&admin)
	db.Create(&user)
	db.Create(&store)
	db.Create(&store1)
	db.Create(&category)
	db.Create(&category1)
	db.Create(&pet)
	db.Create(&pet2)
	db.Create(&product)
	db.Create(&product1)

}
