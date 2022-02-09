package main

import (
	"fmt"
	"petshop/config"
	categoryCtrl "petshop/delivery/controller/category"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	categoryRepo "petshop/repository/category"
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

	categoryRepository := categoryRepo.NewCategoryRepository(db)
	categoryController := categoryCtrl.NewCategoryController(categoryRepository)

	route.RegisterPath(e, categoryController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
