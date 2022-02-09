package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"petshop/config"
	categoryCtrl "petshop/delivery/controller/category"
	productController "petshop/delivery/controller/product"
	transactionController "petshop/delivery/controller/transaction"
	userCtrl "petshop/delivery/controller/user"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	categoryRepo "petshop/repository/category"
	productRepo "petshop/repository/product"
	transactionRepo "petshop/repository/transaction"
	userRepo "petshop/repository/user"
	"petshop/util"
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

	userRepository := userRepo.NewUserRepository(db)
	userController := userCtrl.NewUserController(userRepository)

	categoryRepository := categoryRepo.NewCategoryRepository(db)
	categoryController := categoryCtrl.NewCategoryController(categoryRepository)

	productRepo := productRepo.NewProductRepository(db)
	productController := productController.NewProductController(productRepo)

	transactionRepo := transactionRepo.NewTransactionRepository(db)
	transactionController := transactionController.NewTransactionController(transactionRepo)

	route.RegisterPath(e, userController, productController, categoryController, transactionController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
