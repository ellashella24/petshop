package main

import (
	"fmt"
	"petshop/config"
	productController "petshop/delivery/controller/product"
	transactionController "petshop/delivery/controller/transaction"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	productRepo "petshop/repository/product"
	transactionRepo "petshop/repository/transaction"
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

	productRepo := productRepo.NewProductRepository(db)
	productController := productController.NewProductController(productRepo)

	transactionRepo := transactionRepo.NewTransactionRepository(db)
	transactionController := transactionController.NewTransactionController(transactionRepo)

	route.RegisterPath(e, productController transactionController)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
