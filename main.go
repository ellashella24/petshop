package main

import (
	"fmt"
	"petshop/config"
	storeCtrl "petshop/delivery/controller/store"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	storeRepo "petshop/repository/store"
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

	storeRepository := storeRepo.NewStoreRepository(db)
	storeController := storeCtrl.NewStoreController(storeRepository)

	route.RegisterPath(e, storeController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
