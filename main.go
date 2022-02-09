package main

import (
	"fmt"
	"petshop/config"
	categoryCtrl "petshop/delivery/controller/category"
	cityCtrl "petshop/delivery/controller/city"
	productController "petshop/delivery/controller/product"
	transactionController "petshop/delivery/controller/transaction"
	userCtrl "petshop/delivery/controller/user"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	categoryRepo "petshop/repository/category"
	cityRepo "petshop/repository/city"
	productRepo "petshop/repository/product"
	transactionRepo "petshop/repository/transaction"
	userRepo "petshop/repository/user"
	"petshop/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config := config.GetConfig()
	db := util.InitDB(config)
	e := echo.New()

	middleware.LogMiddleware(e)

	e.Validator = &CustomValidator{validator: validator.New()}

	cityRepository := cityRepo.NewCityRepository(db)
	cityController := cityCtrl.NewCityController(cityRepository)

	userRepository := userRepo.NewUserRepository(db)
	userController := userCtrl.NewUserController(userRepository)

	categoryRepository := categoryRepo.NewCategoryRepository(db)
	categoryController := categoryCtrl.NewCategoryController(categoryRepository)

	productRepo := productRepo.NewProductRepository(db)
	productController := productController.NewProductController(productRepo)

	transactionRepo := transactionRepo.NewTransactionRepository(db)
	transactionController := transactionController.NewTransactionController(transactionRepo)

	route.RegisterPath(e, userController, productController, categoryController, transactionController, cityController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
